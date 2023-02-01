package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charconstpointer/sammy/github"
	"github.com/charconstpointer/sammy/gpt3"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	OpenAIToken string `envconfig:"OPEN_AI_TOKEN" required:"true"`
	GithubToken string `envconfig:"GH_TOKEN" required:"true"`
}

var (
	config     Config
	maxTokens  = flag.Int("max_tokens", 100, "max tokens cost")
	githubUser = flag.String("user", "charconstpointer", "github user")
)

type App struct {
	ghc *github.Client
	gpt *gpt3.Client
}

func NewApp(g github.Client, s gpt3.Client) *App {
	return &App{
		ghc: &g,
		gpt: &s,
	}
}

func (a *App) SummarizeAcitivity(ctx context.Context) (string, error) {
	ev, err := a.ghc.GetEvents(context.Background(), *githubUser)
	if err != nil {
		return "", fmt.Errorf("failed to get events: %w", err)
	}

	var sb strings.Builder
	for _, e := range ev {
		sb.WriteString(e.Body)
	}

	feed := sb.String()
	summary, err := a.gpt.Summarize(context.Background(), feed, gpt3.WithMaxTokens(*maxTokens))
	if err != nil {
		return "", fmt.Errorf("failed to summarize: %w", err)
	}

	return summary, nil
}

func main() {
	flag.Parse()
	envconfig.MustProcess("", &config)
	t := os.Getenv("GH_TOKEN")
	if t == "" {
		log.Fatal("GH_TOKEN env variable is not set")
	}

	ghc := github.NewClient(t)
	gpt := gpt3.NewClient(config.OpenAIToken, gpt3.DaVinci)
	app := NewApp(*ghc, *gpt)

	summary, err := app.SummarizeAcitivity(context.Background())
	if err != nil {
		log.Fatalf("Failed to summarize: %v", err)
	}
	fmt.Println(summary)
}
