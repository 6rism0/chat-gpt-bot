package bot_test

import (
	"os"
	"testing"

	"github.com/6rism0/chat-gpt-bot/internal/bot"
)

var updatePrivate = bot.Update{
	Id: 1,
	Message: bot.Message{
		Text: "@ChatGPTBot What is the capital of Germany?",
		Chat: bot.Chat{
			Id:       2,
			ChatType: bot.Private,
		},
	},
}

func TestGroup(t *testing.T) {
	updateSupergroup := updatePrivate
	updateSupergroup.Message.Chat.ChatType = bot.Supergroup
	succ, err := bot.Sanitize(updateSupergroup.Message)
	if err == nil {
		t.Errorf("sanatize error for %s should not succeed: %+v", updateSupergroup.Message.Chat.ChatType, succ)
	}
}

func TestHelp(t *testing.T) {
	updatePrivate.Message.Text = "/help"
	succ, err := bot.Sanitize(updatePrivate.Message)
	if err != nil {
		t.Errorf("sanatize error for %s should succeed", err.Error())
	}
	t.Logf("\n Success - %+v", succ)
}

func TestImageEmpty(t *testing.T) {
	updatePrivate.Message.Text = "/image"
	_, err := bot.Sanitize(updatePrivate.Message)
	if err == nil {
		t.Errorf("sanatize error for %s should not succeed", err.Error())
	}
}

func TestImageNotEmpty(t *testing.T) {
	updatePrivate.Message.Text = "/image input text"
	succ, err := bot.Sanitize(updatePrivate.Message)
	if err != nil {
		t.Errorf("sanatize error for %s should succeed", err.Error())
	}
	t.Logf("\n Success - %+v", succ)
}

func TestStart(t *testing.T) {
	updatePrivate.Message.Text = "/start input text"
	succ, err := bot.Sanitize(updatePrivate.Message)
	if err != nil {
		t.Errorf("sanatize error for %s should succeed", err.Error())
	}
	t.Logf("\n Success - %+v", succ)
}

func TestEmpty(t *testing.T) {
	updatePrivate.Message.Text = "@ChatGPTBot "
	succ, err := bot.Sanitize(updatePrivate.Message)
	if err == nil {
		t.Errorf("sanatize error nil should not succeed - %+v", succ)
	}
}

func TestSendResponse(t *testing.T) {
	apiToken := os.Getenv("BOT_TOKEN")
	if apiToken == "" {
		t.Skip("Skipping testing against production set Telegram BOT TOKEN to run test")
	}

	succ, err := bot.SendTextToTelegramChat(53333233, "TestSendTextToTelegramChat")
	if err != nil {
		t.Errorf("Sending text response failed - %s", err.Error())
	}
	t.Logf("\n Send response - %s", succ)

	textResponse := bot.TextResponse{
		ChatId: 53333233,
		Text:   "TestSendResponse",
	}
	succ, err = bot.SendResponseToTelegramChat(textResponse)
	if err != nil {
		t.Errorf("Sending text response failed - %s", err.Error())
	}
	t.Logf("\n Send response - %s", succ)

	imageResponse := bot.ImageResponse{
		ChatId:   53333233,
		ImageUrl: "https://images.ctfassets.net/lzny33ho1g45/T5qqQQVznbZaNyxmHybDT/b76e0ff25a495e00647fa9fa6193a3c2/best-url-shorteners-00-hero.png?w=1520&fm=jpg&q=30&fit=thumb&h=760",
	}
	succ, err = bot.SendResponseToTelegramChat(imageResponse)
	if err != nil {
		t.Errorf("Sending image response failed - %s", err.Error())
	}
	t.Logf("\n Send response - %s", succ)
}
