package models

type Country struct {
	Id         int
	Name       string
	Code       string
	Population int
}

func (c Country) GetId() int {
	return c.Id
}
