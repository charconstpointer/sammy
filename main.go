package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
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
	maxTokens  = flag.Int("max_tokens", 200, "max tokens cost")
	githubUser = flag.String("user", "charconstpointer", "github user")
	public     = flag.Bool("public", false, "only public events")
	from       = flag.String("from", "", "from date in format of RFC3339")
	to         = flag.String("to", "", "to date in format of RFC3339")
	verbose    = flag.Bool("verbose", false, "verbose logging")
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
	log.Printf("Making feed from %d events", len(ev))
	var sb strings.Builder
	for _, e := range ev {
		log.Printf("Masking event: %s", e.Body)
		for _, t := range e.Tokens {
			_ = a.masker.Add(t)
		}
		sb.WriteString(e.Body)
	}

	return a.masker.MaskString(sb.String())
}

func mustParseTimeRange(from, to string) (time.Time, time.Time) {
	parseTime := func(x string) time.Time {
		if x == "" {
			return time.Now().UTC()
		}

		t, err := time.Parse(time.RFC3339, x)
		if err != nil {
			panic(err)
		}
		return t
	}

	return parseTime(from), parseTime(to)
}

func main() {
	flag.Parse()
	var c Config
	envconfig.MustProcess("", &c)
	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}

	var (
		ghc    = github.NewClient(c.GithubToken)
		gpt    = gpt3.NewClient(c.OpenAIToken, gpt3.DaVinci)
		masker = masker.New()
		app    = NewApp(ghc, gpt, masker)
		ctx    = context.Background()
	)

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	from, to := mustParseTimeRange(*from, *to)
	if from.Equal(to) {
		to.Add(time.Hour * 24)
	}
	log.Printf("Generating report for date range from %s to %s", from, to)
	summary, err := app.SummarizeAcitivity(ctx, *githubUser, *public, from, to)
	if err != nil {
		log.Fatalf("Failed to summarize: %v", err)
	}
	log.Println(summary)
}
