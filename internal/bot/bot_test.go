package bot_test

import (
	"testing"

	"github.com/6rism0/chat-gpt-bot/internal/bot"
)

var updatePrivate = bot.Update{
	Id: 1,
	Message: bot.Message{
		Text: "Test",
		Chat: bot.Chat{
			Id:       2,
			ChatType: bot.Private,
		},
	},
}

func TestSanatize(t *testing.T) {
	_, err := bot.Sanitize(updatePrivate.Message)
	if err != nil {
		t.Errorf("sanatize error for %s: %s", updatePrivate.Message.Chat.ChatType, err.Error())
	}
	updateGroup := updatePrivate
	updateGroup.Message.Chat.ChatType = bot.Group
	_, err = bot.Sanitize(updateGroup.Message)
	if err != nil {
		t.Errorf("sanatize error for %s: %s", updateGroup.Message.Chat.ChatType, err.Error())
	}
	updateSupergroup := updateGroup
	updateSupergroup.Message.Chat.ChatType = bot.Supergroup
	succ, err := bot.Sanitize(updateSupergroup.Message)
	if err == nil {
		t.Errorf("sanatize error for %s should not succeed: %s", updateSupergroup.Message.Chat.ChatType, succ)
	}
	updateUndefined := updateSupergroup
	updateUndefined.Message.Chat.ChatType = bot.Undefined
	succ, err = bot.Sanitize(updateUndefined.Message)
	if err == nil {
		t.Errorf("sanatize error for %s should not succeed: %s", updateUndefined.Message.Chat.ChatType, succ)
	}
}
