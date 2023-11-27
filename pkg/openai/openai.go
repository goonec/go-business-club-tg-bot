package openai

import (
	"context"
	"errors"
	"github.com/sashabaranov/go-openai"
	"log"
	"sync"
	"time"
)

type OpenAI struct {
	client *openai.Client
	prompt string
	model  string
	mu     sync.Mutex
}

func NewOpenAIConnect(apiKey string) *OpenAI {
	ai := &OpenAI{
		client: openai.NewClient(apiKey),
		model:  openai.GPT3Dot5Turbo,
	}

	if apiKey == "" {
		log.Println("[WARNING] api key is empty")
	}

	return ai
}

func (o *OpenAI) ResponseGPT(text string) (string, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	request := openai.ChatCompletionRequest{
		Model: o.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: text,
			},
		},
		MaxTokens:   800,
		Temperature: 1,
		TopP:        1,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	resp, err := o.client.CreateChatCompletion(ctx, request)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("no choices in openai response")
	}

	return resp.Choices[0].Message.Content, nil
}
