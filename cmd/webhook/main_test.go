package webhook_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/6rism0/chat-gpt-bot/cmd/webhook"
	"github.com/6rism0/chat-gpt-bot/internal/ai"
	"github.com/6rism0/chat-gpt-bot/internal/bot"
	"github.com/6rism0/chat-gpt-bot/internal/image"
	"github.com/6rism0/chat-gpt-bot/internal/util"
)

func TestResponse(t *testing.T) {
	// t.Setenv("OPENAI_API_KEY", "")
	apiToken := os.Getenv("OPENAI_API_KEY")
	if apiToken == "" {
		t.Skip("Skipping testing against production set OPENAI_API_KEY to run test")
	}
	message := bot.Message{
		Text: "What is the capital of Germany?",
		Chat: bot.Chat{
			Id:       123,
			ChatType: bot.Private,
		},
	}

	response, err := webhook.CreateResponse(message, bot.CommandInput{
		Input: "Help Command",
	})
	if err != nil {
		t.Errorf("error - %s", err.Error())
		return
	}
	t.Logf("response - %+v", response)

	response, err = webhook.CreateResponse(message, bot.ImageInput{
		Input: "Help Command",
	})
	if err != nil {
		t.Errorf("error - %s", err.Error())
		return
	}
	t.Logf("response - %+v", response)

	response, err = webhook.CreateResponse(message, bot.TextInput{
		Input: "What is the capital of Germany?",
	})
	if err != nil {
		t.Errorf("error - %s", err.Error())
		return
	}
	t.Logf("response - %+v", response)
}

func TestImage(t *testing.T) {
	// t.Setenv("OPENAI_API_KEY", "")
	apiToken := os.Getenv("OPENAI_API_KEY")
	if apiToken == "" {
		t.Skip("Skipping testing against production set OPENAI_API_KEY to run test")
	}
	request := image.DefaultRequest("fire spitting wiener dog")
	util.LogDebug(fmt.Sprintf("start request for %+v", request))
	client := ai.OpenAIClient(os.Getenv("OPENAI_API_KEY"))
	res, err := image.RequestImage(client, request)
	if err != nil {
		t.Errorf("Could not complete request - %s", err.Error())
		return
	}
	if len(res.Data) > 0 {
		t.Log("\n \n \n")
		for i := 0; i < len(res.Data)-1; i++ {
			t.Logf("Data %d -> %s \n", i, res.Data[i])
		}
	} else {
		t.Error("Empty Data")
	}
}
