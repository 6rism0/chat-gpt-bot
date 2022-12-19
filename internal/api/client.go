package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	e "github.com/6rism0/chat-gpt-bot/internal/error"
)

type Client struct {
	HttpClient http.Client
	BaseURL    string
	Token      string
}

const jsonType = "application/json"

func (client *Client) SendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Accept", jsonType)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.Token))
	req.Header.Set("Content-Type", jsonType)

	res, err := client.HttpClient.Do(req)
	if err != nil {
		fmt.Printf("res error %s", err)
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes e.ErrorResponse
		err = json.NewDecoder(res.Body).Decode(&errRes)
		if err != nil || errRes.Error == nil {
			return fmt.Errorf("error, status code: %d", res.StatusCode)
		}
		return fmt.Errorf("error, status code: %d, message: %s", res.StatusCode, errRes.Error.Message)
	}

	data, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(data, &v)
	if err != nil {
		fmt.Printf("decode error %s", err)
		return err
	}
	return nil
}
