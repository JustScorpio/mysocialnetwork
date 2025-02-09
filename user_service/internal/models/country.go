package models

type User struct {
	Id   int
	Name string
	Mail string
}

func (c User) GetId() int {
	return c.Id
}
