package webhook

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/6rism0/chat-gpt-bot/internal/ai"
	"github.com/6rism0/chat-gpt-bot/internal/bot"
	"github.com/6rism0/chat-gpt-bot/internal/completion"
	"github.com/6rism0/chat-gpt-bot/internal/image"
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
		send(bot.TextResponse{ChatId: update.Message.Chat.Id, Text: err.Error()})
		return
	}

	response, err := CreateResponse(update.Message, sanitizedSeed)
	if err != nil {
		send(bot.TextResponse{ChatId: update.Message.Chat.Id, Text: err.Error()})
		return
	}

	send(response)
}

func CreateResponse(message bot.Message, input bot.Input) (bot.Response, error) {
	client := ai.OpenAIClient(os.Getenv(tokenENV))
	switch input.(type) {
	case bot.CommandInput:
		return bot.TextResponse{
			Text:   input.Text(),
			ChatId: message.Chat.Id,
		}, nil
	case bot.TextInput:
		request := completion.DefaultCompletion(input.Text())
		res, err := completion.RequestCompletion(client, request)
		if err != nil {
			util.LogError(fmt.Sprintf("Could not complete request - %s", err.Error()))
			return nil, err
		}
		if len(res.Choices) > 0 {
			return bot.TextResponse{
				Text:   res.Choices[0].Text,
				ChatId: message.Chat.Id,
			}, nil
		}
	case bot.ImageInput:
		request := image.DefaultRequest(input.Text())
		res, err := image.RequestImage(client, request)
		if err != nil {
			util.LogError(fmt.Sprintf("Could not complete request - %s", err.Error()))
			return nil, err
		}
		if len(res.Data) > 0 {
			return bot.ImageResponse{
				ImageUrl: res.Data[0].Url,
				ChatId:   message.Chat.Id,
			}, nil
		}
	}
	return nil, errors.New("no Response")
}

func send(response bot.Response) {
	responseTelegram, errTelegram := bot.SendResponseToTelegramChat(response)
	if errTelegram != nil {
		util.LogError(fmt.Sprintf("got error %s from telegram, response body is %s", errTelegram.Error(), responseTelegram))
	} else {
		util.LogDebug(fmt.Sprintf("successful send to id %d", response.Id()))
	}
}
