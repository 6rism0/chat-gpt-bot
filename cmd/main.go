package main

import (
	"fmt"
	"os"

	"github.com/6rism0/chat-gpt-bot/internal/ai"
	"github.com/6rism0/chat-gpt-bot/internal/completion"
)

const tokenENV = "OPENAI_API_KEY"

func main() {
	request := completion.DefaultCompletion()
	request.Prompt = "What is the answer of everything?"
	client := ai.OpenAIClient(os.Getenv(tokenENV))
	res, err := completion.RequestCompletion(client, *request)
	if err != nil {
		fmt.Printf("Could not complete request - err: %s", err.Error())
		return
	}
	fmt.Printf("Request completed - response: %s", res.Choices[0].Text)
}
