package chatgptbot

import (
	"net/http"

	"github.com/6rism0/chat-gpt-bot/cmd/webhook"
)

func Run(w http.ResponseWriter, r *http.Request) {
	webhook.HandleTelegramWebHook(w, r)
}
