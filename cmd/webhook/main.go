package webhook

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/6rism0/chat-gpt-bot/internal/ai"
	"github.com/6rism0/chat-gpt-bot/internal/bot"
	"github.com/6rism0/chat-gpt-bot/internal/completion"
	"github.com/6rism0/chat-gpt-bot/internal/util"
)

const tokenENV = "OPENAI_API_KEY"

func HandleTelegramWebHook(w http.ResponseWriter, r *http.Request) {

	var update, err = bot.ParseTelegramRequest(r)
	if err != nil {
		util.LogError(fmt.Sprintf("error parsing update, %s", err.Error()))
		return
	} else {
		util.LogDebug(fmt.Sprintf("successful parsed incomming message %s", update.Message.Text))
	}

	sanitizedSeed, err := bot.Sanitize(update.Message)
	if err != nil {
		send(update.Message.Chat.Id, err.Error())
		return
	}

	response, err := CreateResponse(update.Message, sanitizedSeed)
	if err != nil {
		send(update.Message.Chat.Id, err.Error())
		return
	}

	send(update.Message.Chat.Id, response)
}

func CreateResponse(messsage bot.Message, text string) (string, error) {
	util.LogDebug(fmt.Sprintf("CreateResponse %+v, text: %+v", messsage, text))
	request := completion.DefaultCompletion(text)
	util.LogDebug(fmt.Sprintf("start request for %+v", request))
	client := ai.OpenAIClient(os.Getenv(tokenENV))
	res, err := completion.RequestCompletion(client, request)
	if err != nil {
		util.LogError(fmt.Sprintf("Could not complete request - %s", err.Error()))
		return "", err
	}
	if len(res.Choices) > 0 {
		return res.Choices[0].Text, nil
	}
	return "", errors.New("no response")
}

func send(id int, text string) {
	responseTelegram, errTelegram := bot.SendTextToTelegramChat(id, text)
	if errTelegram != nil {
		util.LogError(fmt.Sprintf("got error %s from telegram, response body is %s", errTelegram.Error(), responseTelegram))
	} else {
		util.LogDebug(fmt.Sprintf("successful send to id %d", id))
	}
}
