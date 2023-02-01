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
		PerPage: 10,
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
		body := fmt.Sprintf("%s has performed an action of %s inside the %s repository\n", e.GetActor().GetLogin(), e.GetType(), e.GetRepo().GetName())
		event := sammy.NewEvent(e.GetActor().GetLogin(), body)
		events = append(events, &event)

	}
	return events, nil
}
