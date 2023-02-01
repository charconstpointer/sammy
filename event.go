package sammy

import "time"

type Event struct {
	CreatedAt time.Time
	Subject   string
	Body      string
}

func NewEvent(subj, body string) Event {
	return Event{
		CreatedAt: time.Now().UTC(),
		Subject:   subj,
		Body:      body,
	}
}
