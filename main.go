package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/charconstpointer/sammy/github"
	"github.com/charconstpointer/sammy/gpt3"
	"github.com/charconstpointer/sammy/namesgenerator"
	"github.com/kelseyhightower/envconfig"
)

var (
	config     Config
	maxTokens  = flag.Int("max_tokens", 100, "max tokens cost")
	githubUser = flag.String("user", "charconstpointer", "github user")
	public     = flag.Bool("public", true, "only public events")
)

type Config struct {
	OpenAIToken string `envconfig:"OPEN_AI_TOKEN" required:"true"`
	GithubToken string `envconfig:"GITHUB_TOKEN" required:"true"`
}

type App struct {
	ghc    *github.Client
	gpt    *gpt3.Client
	masker *namesgenerator.Masker
}

func NewApp(g *github.Client, s *gpt3.Client, masker *namesgenerator.Masker) *App {
	return &App{
		ghc:    g,
		gpt:    s,
		masker: masker,
	}
}

func (a *App) SummarizeAcitivity(ctx context.Context, user string, public bool) (string, error) {
	ev, err := a.ghc.UserEvents(ctx, user, public, time.Now().Add(-time.Hour*1), time.Now())
	if err != nil {
		return "", fmt.Errorf("failed to get events: %w", err)
	}

	var sb strings.Builder
	for _, e := range ev {
		for _, t := range e.Tokens {
			a.masker.Register(t)
		}
		sb.WriteString(e.Body)
	}

	feed := sb.String()
	feed = a.masker.MaskString(feed)
	summary, err := a.gpt.Summarize(ctx, feed, gpt3.WithMaxTokens(*maxTokens))
	if err != nil {
		return "", fmt.Errorf("failed to summarize: %w", err)
	}

	return a.masker.UnmaskString(summary), nil
}

func main() {
	flag.Parse()
	envconfig.MustProcess("", &config)

	ghc := github.NewClient(config.GithubToken)
	gpt := gpt3.NewClient(config.OpenAIToken, gpt3.DaVinci)
	masker := namesgenerator.NewMasker()
	app := NewApp(ghc, gpt, masker)

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	summary, err := app.SummarizeAcitivity(ctx, *githubUser, *public)
	if err != nil {
		log.Fatalf("Failed to summarize: %v", err)
	}
	fmt.Println(summary)
}
