package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/charconstpointer/sammy/github"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	OpenAIToken string `envconfig:"OPEN_AI_TOKEN" required:"true"`
	GithubToken string `envconfig:"GH_TOKEN" required:"true"`
}

var (
	config     Config
	sample     string
	sampleSrc  = flag.String("sample", "sample.txt", "sample text file")
	maxTokens  = flag.Int("max_tokens", 100, "max tokens")
	temp       = flag.Float64("temp", 0.3, "temperature")
	githubUser = flag.String("user", "charconstpointer", "github user")
)

func main() {
	flag.Parse()
	envconfig.MustProcess("", &config)
	t := os.Getenv("GH_TOKEN")
	if t == "" {
		log.Fatal("GH_TOKEN env variable is not set")
	}
	c := github.NewClient(t)
	ev, err := c.GetEvents(context.Background(), *githubUser)
	if err != nil {
		log.Fatal(err)
	}
	var sb strings.Builder
	for _, e := range ev {
		sb.WriteString(e.Body)
	}
	feed := sb.String()
	summary, err := summarize(context.Background(), feed)
	if err != nil {
		log.Fatalf("Failed to summarize: %v", err)
	}
	log.Printf("Summary: %s", summary)
}

func summarize(ctx context.Context, feed string) (string, error) {
	prompt := fmt.Sprintf("Pick releveant information about repositories and actions made from the following text, keep in mind actions should be reporter in a human readable form: %s", feed)

	c := gpt3.NewClient(config.OpenAIToken, gpt3.WithDefaultEngine("text-davinci-003"))
	log.Printf("Requesting summary of: %s", prompt)
	req := gpt3.CompletionRequest{
		Prompt:    []string{prompt},
		MaxTokens: maxTokens,
		Echo:      false,
	}
	res, err := c.Completion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to get completion: %w", err)
	}
	return res.Choices[0].Text, nil
}
