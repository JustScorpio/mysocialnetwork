package models

import "time"

type Chat struct {
	Id           int
	Owner        User
	CreatedAt    time.Time
	Messages     []Message
	Participants []User
}

func (c Chat) GetId() int {
	return c.Id
}
