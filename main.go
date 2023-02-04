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
	"github.com/charconstpointer/sammy/masker"
	"github.com/kelseyhightower/envconfig"
)

var (
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
	masker *masker.Masker
}

func NewApp(g *github.Client, s *gpt3.Client, masker *masker.Masker) *App {
	return &App{
		ghc:    g,
		gpt:    s,
		masker: masker,
	}
}

func (a *App) SummarizeAcitivity(ctx context.Context, user string, public bool, start, end time.Time) (string, error) {
	ev, err := a.ghc.UserEvents(ctx, user, public, start, end)
	if err != nil {
		return "", fmt.Errorf("failed to get events: %w", err)
	}

	feed := a.makeFeed(ev)
	summary, err := a.gpt.Summarize(ctx, feed, gpt3.WithMaxTokens(*maxTokens))
	if err != nil {
		return "", fmt.Errorf("failed to summarize: %w", err)
	}

	return a.masker.UnmaskString(summary), nil
}

func (a *App) makeFeed(ev []*github.Event) string {
	var sb strings.Builder
	for _, e := range ev {
		for _, t := range e.Tokens {
			a.masker.MustAdd(t)
		}
		sb.WriteString(e.Body)
	}

	return a.masker.MaskString(sb.String())
}

func main() {
	flag.Parse()
	var c Config
	envconfig.MustProcess("", &c)

	var (
		ghc    = github.NewClient(c.GithubToken)
		gpt    = gpt3.NewClient(c.OpenAIToken, gpt3.DaVinci)
		masker = masker.New()
		app    = NewApp(ghc, gpt, masker)
		ctx    = context.Background()
	)

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	from, to := time.Now().UTC().Add(-time.Hour*24), time.Now().UTC()
	summary, err := app.SummarizeAcitivity(ctx, *githubUser, *public, from, to)
	if err != nil {
		log.Fatalf("Failed to summarize: %v", err)
	}
	log.Println(summary)
}
