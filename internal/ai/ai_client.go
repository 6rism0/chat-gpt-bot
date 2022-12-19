package ai

import (
	"net/http"

	api "github.com/6rism0/chat-gpt-bot/internal/api"
)

const openAIUrl = "https://api.openai.com/v1"

func OpenAIClient(authToken string) *api.Client {
	return &api.Client{
		HttpClient: *http.DefaultClient,
		BaseURL:    openAIUrl,
		Token:      authToken,
	}
}
