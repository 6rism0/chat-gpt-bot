package chatgptbot

import (
	"net/http"

	"github.com/6rism0/chat-gpt-bot/cmd/bot"
)

func Run(w http.ResponseWriter, r *http.Request) {
	bot.HandleTelegramWebHook(w, r)
}
