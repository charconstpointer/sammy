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

func (c *Client) GetEvents(ctx context.Context, user string) ([]*sammy.Event, error) {
	ev, res, err := c.c.Activity.ListEventsPerformedByUser(ctx, user, false, &github.ListOptions{
		PerPage: 25,
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
		p, err := e.ParsePayload()
		if err != nil {
			return nil, fmt.Errorf("failed to parse payload: %w", err)
		}
		var body string
		switch *e.Type {
		case "IssuesEvent":
			body = c.handleIssue(ctx, p.(*github.IssuesEvent))
		case "IssueCommentEvent":
			body = c.handleIssueComment(ctx, p.(*github.IssueCommentEvent))
		case "PullRequestEvent":
			body = c.handlePullRequest(ctx, p.(*github.PullRequestEvent))
		case "PullRequestReviewEvent":
			body = c.handlePullRequestReview(ctx, p.(*github.PullRequestReviewEvent))
		case "PullRequestReviewCommentEvent":
			body = c.handlePullRequestReviewComment(ctx, p.(*github.PullRequestReviewCommentEvent))
		case "PushEvent":
			body = c.handlePush(ctx, p.(*github.PushEvent))
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

func (c *Client) handleIssueComment(ctx context.Context, e *github.IssueCommentEvent) string {
	var (
		issueTitle = "an issue"
	)

	if e.Issue != nil {
		issueTitle = *e.Issue.Title
	}

	return fmt.Sprintf("commented on an issue %s", issueTitle)
}

func (c *Client) handlePullRequest(ctx context.Context, e *github.PullRequestEvent) string {
	var (
		pullRequestTitle, pullRequestURL = "a pull request", "a pull request"
	)
	if e.PullRequest != nil {
		pullRequestTitle = *e.PullRequest.Title
	}
	if e.Repo != nil {
		pullRequestURL = *e.Repo.PullsURL
	}

	return fmt.Sprintf("created a pull request %s avaiable at %s", pullRequestTitle, pullRequestURL)
}

func (c *Client) handlePullRequestReview(ctx context.Context, e *github.PullRequestReviewEvent) string {
	var (
		pullRequestTitle, pullRequestURL = "a pull request", "a pull request"
	)
	if e.PullRequest != nil {
		pullRequestTitle = *e.PullRequest.Title
	}
	if e.Repo != nil {
		pullRequestURL = *e.Repo.PullsURL
	}

	return fmt.Sprintf("created a pull request %s available at %s", pullRequestTitle, pullRequestURL)
}

func (c *Client) handlePullRequestReviewComment(ctx context.Context, e *github.PullRequestReviewCommentEvent) string {
	var (
		pullRequestTitle = "a pull request"
	)
	if e.PullRequest != nil {
		pullRequestTitle = *e.PullRequest.Title
	}
	return fmt.Sprintf("commented on a pull request %s", pullRequestTitle)
}

func (c *Client) handlePush(ctx context.Context, e *github.PushEvent) string {
	var (
		repoName = "a commit"
	)
	if e.Repo != nil {
		repoName = *e.Repo.PullsURL
	}
	return fmt.Sprintf("pushed to %s", repoName)
}
