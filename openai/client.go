package client

import (
	"context"
	"gollm/memory"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/sashabaranov/go-openai"
)

type Client struct {
	client *openai.Client
	once   sync.Once
	apiKey string
	Memory *memory.MemoryManager
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		Memory: memory.NewMemoryManager(),
	}
}

func (c *Client) GetClient() *openai.Client {
	c.once.Do(func() {
		config := openai.DefaultConfig(c.apiKey)
		config.BaseURL = "https://api.openai.com/v1"
		c.client = openai.NewClientWithConfig(config)
		// log.Info("openai api client initialized", "apiKey", c.apiKey)
	})
	return c.client
}

func (c *Client) CreateChatCompletion(ctx context.Context, request openai.ChatCompletionRequest) (string, error) {
	client := c.GetClient()

	c.Memory.AddMessage(request.Messages[0])
	// log.Info("[client] user message added", "message", request.Messages[0])

	history := c.Memory.GetHistory()
	// log.Info("[client] message history retrieved", "history", history)

	request.Messages = history
	// log.Info("[client] request message history set", "request.Messages", request.Messages)

	resp, err := client.CreateChatCompletion(ctx, request)
	if err != nil {
		log.Error("[client] could not make chat completion request", "err", err)
		return "", err
	}

	// log.Info("[client] chat completion created by client", "completion", resp.Choices[0].Message.Content)

	c.Memory.AddMessage(resp.Choices[0].Message)
	// log.Info("[client] chat completion message added to history")

	return resp.Choices[0].Message.Content, nil
}
