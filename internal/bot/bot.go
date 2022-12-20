package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"

	"github.com/6rism0/chat-gpt-bot/internal/util"
)

const telegramApiBaseUrl string = "https://api.telegram.org/bot"
const telegramApiSendMessage string = "/sendMessage"
const telegramTokenEnv string = "BOT_TOKEN"

const cleanTextRegex = `(@\S+)`

var telegramApi string = telegramApiBaseUrl + os.Getenv(telegramTokenEnv) + telegramApiSendMessage

type Update struct {
	Id      int     `json:"update_id"`
	Message Message `json:"message"`
}

type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}

type Chat struct {
	Id       int  `json:"id"`
	ChatType Type `json:"type"`
}

type Type string

const (
	Undefined  Type = ""
	Private    Type = "private"
	Group      Type = "group"
	Supergroup Type = "supergroup"
)

func ParseTelegramRequest(r *http.Request) (*Update, error) {
	var update Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		util.LogError(fmt.Sprintf("could not decode incoming update %s", err.Error()))
		return nil, err
	}
	return &update, nil
}

func SendTextToTelegramChat(chatId int, text string) (string, error) {
	response, err := http.PostForm(
		telegramApi,
		url.Values{
			"chat_id": {strconv.Itoa(chatId)},
			"text":    {text},
		})
	if err != nil {
		util.LogError(fmt.Sprintf("error when posting text to chat - %s", err.Error()))
		return "", err
	}
	defer response.Body.Close()

	bodyBytes, errRead := io.ReadAll(response.Body)
	if errRead != nil {
		util.LogError(fmt.Sprintf("error parsing telegram answer %s", err.Error()))
		return "", errRead
	}

	return string(bodyBytes), nil
}

func Sanitize(msg Message) (string, error) {
	var err error
	var message string
	switch msg.Chat.ChatType {
	case Private:
		message, err = Strip(msg.Text)
	case Group:
		message, err = Strip(msg.Text)
	case Undefined:
		message = ""
		err = errors.New("can't determine chat type")
	default:
		message = ""
		err = errors.New("bot can only be used in group or private chat")
	}
	//util.LogDebug(fmt.Sprintf("Sanitize message: %s -> msg %s, err %s", msg.Text, message, err.Error()))
	return message, err
}

func Strip(text string) (string, error) {
	r, err := regexp.Compile(cleanTextRegex)
	if err != nil {
		return "", err
	}
	var strippedText = r.ReplaceAllString(text, "")
	if strippedText != "" {
		return strippedText, nil
	} else {
		return "", errors.New("empty message")
	}
}
