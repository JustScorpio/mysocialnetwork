package models

type User struct {
	Id       int
	UserName string
	Mail     string
}

func (c User) GetId() int {
	return c.Id
}
