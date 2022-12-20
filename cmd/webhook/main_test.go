package webhook_test

// import (
// 	"testing"

// 	"github.com/6rism0/chat-gpt-bot/cmd/webhook"
// 	"github.com/6rism0/chat-gpt-bot/internal/bot"
// )

// func TestResponse(t *testing.T) {
// 	t.Setenv("OPENAI_API_KEY", "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")
// 	message := bot.Message{
// 		Text: "What is the capital of Germany?",
// 		Chat: bot.Chat{
// 			Id:       123,
// 			ChatType: bot.Private,
// 		},
// 	}

// 	response, err := webhook.CreateResponse(message, "What is the capital of Germany?")
// 	if err != nil {
// 		t.Errorf("error - %s", err.Error())
// 		return
// 	}
// 	t.Logf("response - %+v", response)
// }
