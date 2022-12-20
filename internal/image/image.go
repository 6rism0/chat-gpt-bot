package image

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/6rism0/chat-gpt-bot/internal/api"
)

const endpoint = "/images/generations"

type ImageRequest struct {
	Prompt string    `json:"prompt,omitempty"`
	N      int       `json:"n,omitempty"`
	Size   ImageSize `json:"size,omitempty"`
}

type ImageResponse struct {
	Data []ImageData `json:"data,omitempty"`
}

type ImageData struct {
	Url string `json:"url,omitempty"`
}

type ImageSize string

const (
	Small  ImageSize = "256x256"
	Medium ImageSize = "512x512"
	Large  ImageSize = "1024x1024"
)

func DefaultRequest(prompt string) ImageRequest {
	return ImageRequest{
		Prompt: prompt,
		N:      2,
		Size:   Medium,
	}
}

func RequestImage(client *api.Client, request ImageRequest) (response ImageResponse, err error) {
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
