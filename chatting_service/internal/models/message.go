package models

import "time"

type Message struct {
	Id       int
	Author   User
	SendTime time.Time
	Content  string
	Chat     Chat
}

func (c Message) GetId() int {
	return c.Id
}
