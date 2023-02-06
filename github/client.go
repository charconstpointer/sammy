package github

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

type Client struct {
	c *github.Client
}

func NewClient(token string) *Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	c := oauth2.NewClient(context.TODO(), ts)
	ghc := github.NewClient(c)
	return &Client{
		c: ghc,
	}
}

type Event struct {
	CreatedAt time.Time
	Body      string
	Tokens    []string
}

func NewEvent(subj, body string, tokens []string) Event {
	return Event{
		CreatedAt: time.Now().UTC(),
		Body:      body,
		Tokens:    tokens,
	}
}

func (c *Client) UserEvents(ctx context.Context, user string, public bool, start, end time.Time) (events []*Event, err error) {
	ev, err := c.githubEvents(ctx, user, public, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	for _, e := range ev {
		var (
			body   string
			tokens []string
		)
		switch *e.Type {
		case "IssuesEvent":
			body, tokens = c.handleIssue(ctx, e)
		case "IssueCommentEvent":
			body, tokens = c.handleIssueComment(ctx, e)
		case "PullRequestEvent":
			body, tokens = c.handlePullRequest(ctx, e)
		case "PullRequestReviewEvent":
			body, tokens = c.handlePullRequestReview(ctx, e)
		case "PullRequestReviewCommentEvent":
			body, tokens = c.handlePullRequestReviewComment(ctx, e)
		case "PushEvent":
			body, tokens = c.handlePush(ctx, e)
		case "WatchEvent":
			body, tokens = c.handleWatch(ctx, e)
		case "ForkEvent":
			body, tokens = c.handleFork(ctx, e)
		default:
			body, tokens = fmt.Sprintf("unhandled event %s", *e.Type), nil
		}
		body = fmt.Sprintf("%s has %s\n", e.GetActor().GetLogin(), body)
		event := NewEvent(e.GetActor().GetLogin(), body, tokens)
		events = append(events, &event)
	}
	return events, nil
}

func (c *Client) githubEvents(ctx context.Context, user string, public bool, start, end time.Time) ([]*github.Event, error) {
	var (
		opt            = &github.ListOptions{PerPage: 10}
		filteredEvents []*github.Event
	)
	for {
		events, resp, err := c.c.Activity.ListEventsPerformedByUser(ctx, user, public, opt)
		if err != nil {
			return nil, err
		}
		for i := len(events) - 1; i >= 0; i-- {
			event := events[i]
			if event.CreatedAt.After(end) {
				continue
			}
			if event.CreatedAt.Before(start) {
				return filteredEvents, nil
			}
			if event.CreatedAt.After(start) && event.CreatedAt.Before(end) {
				filteredEvents = append(filteredEvents, event)
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return filteredEvents, nil
}

func (c *Client) handleIssue(ctx context.Context, e *github.Event) (string, []string) {
	var (
		p               = e.Payload()
		ev              = p.(*github.IssuesEvent)
		issueTitle      = *ev.Issue.Title
		issueURL        = *ev.Issue.URL
		sensitiveTokens = []string{issueTitle, issueURL}
	)

	return fmt.Sprintf("created %s issue available at %s", issueTitle, issueURL), sensitiveTokens
}

func (c *Client) handleIssueComment(ctx context.Context, e *github.Event) (string, []string) {
	var (
		p               = e.Payload()
		ev              = p.(*github.IssueCommentEvent)
		issueTitle      = *ev.Issue.Title
		issueURL        = *ev.Issue.URL
		sensitiveTokens = []string{issueTitle, issueURL}
	)
	return fmt.Sprintf("commented on an issue %s", issueTitle), sensitiveTokens
}

func (c *Client) handlePullRequest(ctx context.Context, e *github.Event) (string, []string) {
	var (
		p                = e.Payload()
		ev               = p.(*github.PullRequestEvent)
		pullRequestTitle = *ev.PullRequest.Title
		pullRequestURL   = *ev.PullRequest.URL
		sensitiveTokens  = []string{pullRequestTitle, pullRequestURL}
	)

	return fmt.Sprintf("created a pull request %s avaiable at %s", pullRequestTitle, pullRequestURL), sensitiveTokens
}

func (c *Client) handlePullRequestReview(ctx context.Context, e *github.Event) (string, []string) {
	var (
		p                = e.Payload()
		ev               = p.(*github.PullRequestReviewEvent)
		pullRequestTitle = *ev.PullRequest.Title
		pullRequestURL   = *ev.PullRequest.URL
		sensitiveTokens  = []string{pullRequestTitle, pullRequestURL}
	)
	return fmt.Sprintf("created a pull request %s available at %s", pullRequestTitle, pullRequestURL), sensitiveTokens
}

func (c *Client) handlePullRequestReviewComment(ctx context.Context, e *github.Event) (string, []string) {
	var (
		p                = e.Payload()
		ev               = p.(*github.PullRequestReviewCommentEvent)
		pullRequestTitle = *ev.PullRequest.Title
		sensitiveTokens  = []string{pullRequestTitle}
	)
	return fmt.Sprintf("commented on a pull request %s", pullRequestTitle), sensitiveTokens
}

func (c *Client) handlePush(ctx context.Context, e *github.Event) (string, []string) {
	var (
		repoName        = *e.Repo.Name
		sensitiveTokens = []string{repoName}
	)
	return fmt.Sprintf("pushed to %s", repoName), sensitiveTokens
}

func (c *Client) handleWatch(ctx context.Context, e *github.Event) (string, []string) {
	var (
		repoName        = *e.Repo.Name
		sensitiveTokens = []string{repoName}
	)
	return fmt.Sprintf("starred the repo %s", repoName), sensitiveTokens
}

func (c *Client) handleFork(ctx context.Context, e *github.Event) (string, []string) {
	var (
		repoName        = *e.Repo.Name
		sensitiveTokens = []string{repoName}
	)
	return fmt.Sprintf("forked the repo %s", repoName), sensitiveTokens
}
