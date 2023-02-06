package gpt3

import (
	"context"
	"fmt"

	"github.com/PullRequestInc/go-gpt3"
)

const prompt = "Based on the following feed of user's github activity, generate a summary of the user's activity, the report should be targeted for higher management person: %s"
const genericPrompth = "Based on the following feed of user's github activity, generate a summary of the user's activity, the report should be targeted for higher management person, you can replace certain actions with more descriptive wording like 'worked on': %s"

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
	prompt := fmt.Sprintf(genericPrompth, text)
	req := gpt3.CompletionRequest{
		Prompt: []string{prompt},
		Echo:   false,
	}
	for _, opt := range opts {
		opt(&req)
	}

	res, err := c.c.Completion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to get completion: %w", err)
	}

	if len(res.Choices) == 0 {
		return "", fmt.Errorf("no choices returned")
	}

	return res.Choices[0].Text, nil
}
