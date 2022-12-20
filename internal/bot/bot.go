package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/6rism0/chat-gpt-bot/internal/util"
)

const telegramApiBaseUrl string = "https://api.telegram.org/bot"
const telegramApiSendMessage string = "/sendMessage"
const telegramApiSendImage string = "/sendPhoto"
const telegramTokenEnv string = "BOT_TOKEN"

const cleanTextRegex = `(@\S+)`

var telegramApi string = telegramApiBaseUrl + os.Getenv(telegramTokenEnv)

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

type Response interface {
	Id() int
	Data() string
}

type ImageResponse struct {
	ChatId   int
	ImageUrl string
}

func (ir ImageResponse) Id() int {
	return ir.ChatId
}

func (ir ImageResponse) Data() string {
	return ir.ImageUrl
}

type TextResponse struct {
	ChatId int
	Text   string
}

func (tr TextResponse) Id() int {
	return tr.ChatId
}

func (tr TextResponse) Data() string {
	return tr.Text
}

type Input interface {
	Text() string
}

type ImageInput struct {
	Input string
}

type TextInput struct {
	Input string
}

type CommandInput struct {
	Input string
}

func (ti TextInput) Text() string {
	return ti.Input
}

func (ii ImageInput) Text() string {
	return ii.Input
}

func (ci CommandInput) Text() string {
	return ci.Input
}

const (
	Start string = "/start"
	Help  string = "/help"
	Image string = "/image"
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

func SendResponseToTelegramChat(r Response) (string, error) {
	var values url.Values
	var sendUrl string
	switch r.(type) {
	case ImageResponse:
		sendUrl = telegramApiBaseUrl + os.Getenv(telegramTokenEnv) + telegramApiSendImage
		values = url.Values{
			"chat_id": {strconv.Itoa(r.Id())},
			"photo":   {r.Data()},
		}
	case TextResponse:
		sendUrl = telegramApi + os.Getenv(telegramTokenEnv) + telegramApiSendMessage
		values = url.Values{
			"chat_id": {strconv.Itoa(r.Id())},
			"text":    {r.Data()},
		}
	default:
		return "", errors.New("unknown type")
	}
	response, err := http.PostForm(sendUrl, values)
	log.Printf("Response - %s", sendUrl)
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

func Sanitize(msg Message) (Input, error) {
	input, err := CheckMessageType(msg)
	if err != nil {
		return nil, err
	}
	return InputToResponse(StripBotName(input))
}

func CheckMessageType(msg Message) (string, error) {
	var err error
	var message string
	switch msg.Chat.ChatType {
	case Private:
		message = msg.Text
		err = nil
	case Group:
		message = msg.Text
		err = nil
	case Undefined:
		message = ""
		err = errors.New("can't determine chat type")
	default:
		message = ""
		err = errors.New("bot can only be used in group or private chat")
	}
	return message, err
}

func StripBotName(text string) string {
	r := regexp.MustCompile(cleanTextRegex)
	return r.ReplaceAllString(text, "")
}

func InputToResponse(inputText string) (Input, error) {
	var input Input
	if strings.Contains(inputText, Start) {
		inputText = "This is a Start Command"
		input = CommandInput{
			Input: inputText,
		}
	} else if strings.Contains(inputText, Help) {
		inputText = "This is a Help Command"
		input = CommandInput{
			Input: inputText,
		}
	} else if strings.Contains(inputText, Image) {
		inputText = strings.ReplaceAll(inputText, Image, "")
		if strings.TrimSpace(inputText) == "" {
			return nil, errors.New("empty input")
		}
		input = ImageInput{
			Input: inputText,
		}
	} else {
		if strings.TrimSpace(inputText) == "" {
			return nil, errors.New("empty input")
		}
		input = TextInput{
			Input: inputText,
		}
	}
	return input, nil
}
