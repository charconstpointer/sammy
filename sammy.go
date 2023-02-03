package sammy

import (
	"time"
)

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

type Masker interface {
	Register(s string)
	MaskString(s string) string
	UnmaskString(s string) string
}
