package completion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/6rism0/chat-gpt-bot/internal/api"
)

type ModelType string

const endpoint = "/completions"

const (
	DavinciText ModelType = "text-davinci-003"
)

type CompletionRequest struct {
	Model     ModelType `json:"model"`
	Prompt    string    `json:"prompt,omitempty"`
	MaxTokens int       `json:"max_tokens,omitempty"`
	Echo      bool      `json:"echo,omitempty"`
}

type CompletionResponse struct {
	ID      string             `json:"id"`
	Object  string             `json:"object"`
	Created uint64             `json:"created"`
	Model   string             `json:"model"`
	Choices []CompletionChoice `json:"choices"`
	Usage   Usage              `json:"usage"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type CompletionChoice struct {
	Text         string `json:"text"`
	Index        int    `json:"index"`
	FinishReason string `json:"finish_reason"`
}

func DefaultCompletion(prompt string) CompletionRequest {
	return CompletionRequest{
		Model:     DavinciText,
		MaxTokens: 100,
		Prompt:    prompt,
		Echo:      true,
	}
}

func RequestCompletion(client *api.Client, request CompletionRequest) (response CompletionResponse, err error) {
	reBytes, err := json.Marshal(request)
	if err != nil {
		fmt.Printf("reBytes error %s", err)
		return
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", client.BaseURL, endpoint), bytes.NewBuffer(reBytes))
	if err != nil {
		fmt.Printf("req error %s", err)
		return
	}
	err = client.SendRequest(req, &response)
	return
}
