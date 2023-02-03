package github

import (
	"context"
	"fmt"
	"net/http"

	"github.com/charconstpointer/sammy"
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

func (c *Client) GetEvents(ctx context.Context, user string, public bool) ([]*sammy.Event, error) {
	ev, res, err := c.c.Activity.ListEventsPerformedByUser(ctx, user, false, &github.ListOptions{
		PerPage: 500,
		Page:    1,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get events: %s", res.Status)
	}
	var events []*sammy.Event
	for _, e := range ev {
		if e.Public == nil || *e.Public != public {
			continue
		}
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
		event := sammy.NewEvent(e.GetActor().GetLogin(), body, tokens)
		events = append(events, &event)
	}
	return events, nil
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
