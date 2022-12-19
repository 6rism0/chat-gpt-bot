package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/6rism0/chat-gpt-bot/internal/util"
)

const telegramApiBaseUrl string = "https://api.telegram.org/bot"
const telegramApiSendMessage string = "/sendMessage"
const telegramTokenEnv string = "BOT_TOKEN"

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

func CreateResponse(msg Message, s string) string {
	var response string
	switch msg.Chat.ChatType {
	case Undefined:
		response = "Error: can't determine chat type"
	case Private:
	case Group:
		response = "TODO"
	default:
		response = "Bot can only be used in group or private chat"
	}
	return response
}

func Sanitize(msg Message) (string, error) {
	switch msg.Chat.ChatType {
	case Private:
	case Group:
		return msg.Text, nil
	case Undefined:
		return "", errors.New("can't determine chat type")
	default:
		return "", errors.New("bot can only be used in group or private chat")
	}
	return "", errors.New("missing type")
}
