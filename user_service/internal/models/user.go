package models

type User struct {
	Id           int
	UserName     string
	Name         string
	PasswordHash string
	Mail         string
	Country      *Country
}

func (c User) GetId() int {
	return c.Id
}
