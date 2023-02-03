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
		p, err := e.ParsePayload()
		if err != nil {
			return nil, fmt.Errorf("failed to parse payload: %w", err)
		}
		var body string
		switch *e.Type {
		case "IssuesEvent":
			body = c.handleIssue(ctx, p.(*github.IssuesEvent))
		case "IssueCommentEvent":
			body = c.handleIssueComment(ctx, e)
		case "PullRequestEvent":
			body = c.handlePullRequest(ctx, e)
		case "PullRequestReviewEvent":
			body = c.handlePullRequestReview(ctx, e)
		case "PullRequestReviewCommentEvent":
			body = c.handlePullRequestReviewComment(ctx, e)
		case "PushEvent":
			body = c.handlePush(ctx, e)
		case "WatchEvent":
			body = c.handleWatch(ctx, e)
		case "ForkEvent":
			body = c.handleFork(ctx, e)
		default:
			body = fmt.Sprintf("unhandled event %s", *e.Type)
		}
		body = fmt.Sprintf("%s has %s\n", e.GetActor().GetLogin(), body)
		event := sammy.NewEvent(e.GetActor().GetLogin(), body)
		events = append(events, &event)
	}
	return events, nil
}

func (c *Client) handleIssue(ctx context.Context, e *github.IssuesEvent) string {
	var (
		issueTitle, issueURL = "an issue", "an issue"
	)

	if e.Issue != nil {
		issueTitle = *e.Issue.Title
	}
	if e.Repo != nil {
		issueURL = *e.Repo.PullsURL
	}

	return fmt.Sprintf("created %s issue available at %s", issueTitle, issueURL)
}

func (c *Client) handleIssueComment(ctx context.Context, e *github.Event) string {
	p := e.Payload()
	ev := p.(*github.IssueCommentEvent)
	return fmt.Sprintf("commented on an issue %s", *ev.Issue.Title)
}

func (c *Client) handlePullRequest(ctx context.Context, e *github.Event) string {
	p := e.Payload()
	ev := p.(*github.PullRequestEvent)

	return fmt.Sprintf("created a pull request %s avaiable at %s", *ev.PullRequest.Title, ev.PullRequest.GetHTMLURL())
}

func (c *Client) handlePullRequestReview(ctx context.Context, e *github.Event) string {
	p := e.Payload()
	ev := p.(*github.PullRequestReviewEvent)
	return fmt.Sprintf("created a pull request %s available at %s", *ev.PullRequest.Title, *ev.PullRequest.URL)
}

func (c *Client) handlePullRequestReviewComment(ctx context.Context, e *github.Event) string {
	p := e.Payload()
	ev := p.(*github.PullRequestReviewCommentEvent)
	return fmt.Sprintf("commented on a pull request %s", *ev.PullRequest.Title)
}

func (c *Client) handlePush(ctx context.Context, e *github.Event) string {
	return fmt.Sprintf("pushed to %s", *e.Repo.Name)
}

func (c *Client) handleWatch(ctx context.Context, e *github.Event) string {
	return fmt.Sprintf("starred the repo %s", *e.Repo.Name)
}

func (c *Client) handleFork(ctx context.Context, e *github.Event) string {
	return fmt.Sprintf("forked the repo %s", *e.Repo.Name)
}
