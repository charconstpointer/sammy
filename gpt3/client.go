package gpt3

import (
	"context"
	"fmt"
	"log"

	"github.com/PullRequestInc/go-gpt3"
)

const prompt = "Pick releveant information about repositories and actions made from the following text, keep in mind actions should be reporter in a human readable form: %s"

type EngineType string

const (
	DaVinci EngineType = "text-davinci-003"
)

type Client struct {
	c gpt3.Client
}

func NewClient(apiToken string, engine EngineType) *Client {
	c := gpt3.NewClient(apiToken, gpt3.WithDefaultEngine(string(engine)))
	return &Client{
		c: c,
	}
}

type CompletionRequestOpt func(*gpt3.CompletionRequest)

func WithMaxTokens(max int) CompletionRequestOpt {
	return func(req *gpt3.CompletionRequest) {
		req.MaxTokens = gpt3.IntPtr(max)
	}
}

func (c *Client) Summarize(ctx context.Context, text string, opts ...CompletionRequestOpt) (string, error) {
	prompt := fmt.Sprintf(prompt, text)
	log.Printf("Prompt: %s\n", prompt)
	req := gpt3.CompletionRequest{
		Prompt: []string{prompt},
		Echo:   false,
	}
	for _, opt := range opts {
		opt(&req)
	}

	log.Printf("Requesting summary of: %s", prompt)
	res, err := c.c.Completion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to get completion: %w", err)
	}

	if len(res.Choices) == 0 {
		return "", fmt.Errorf("no choices returned")
	}

	return res.Choices[0].Text, nil
}
