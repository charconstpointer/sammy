package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	APIKey string `envconfig:"API_KEY" required:"true"`
}

var (
	config    Config
	sample    string
	sampleSrc = flag.String("sample", "sample.txt", "sample text file")
	maxTokens = flag.Int("max_tokens", 100, "max tokens")
	temp      = flag.Float64("temp", 0.3, "temperature")
)

func main() {
	flag.Parse()
	envconfig.MustProcess("", &config)
	c := gpt3.NewClient(config.APIKey)
	s, err := getSample()
	if err != nil {
		log.Fatalf("Failed to get sample: %v", err)
	}
	req := gpt3.CompletionRequest{
		Prompt:           []string{s},
		MaxTokens:        maxTokens,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		TopP:             gpt3.Float32Ptr(1),
		Temperature:      gpt3.Float32Ptr(float32(*temp)),
	}
	res, err := c.Completion(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to get completion: %v", err)
	}
	log.Printf("Summary: %v", res.Choices[0].Text)
}

func getSample() (string, error) {
	c, err := os.ReadFile(*sampleSrc)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Summarize text of user activity on github below in a non JSON, human readable format\n%s", c), nil
}
