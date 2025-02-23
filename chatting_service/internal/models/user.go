package models

type User struct {
	Id       int
	UserName string
	Name     string
	Chats    []Chat
}

func (c User) GetId() int {
	return c.Id
}
