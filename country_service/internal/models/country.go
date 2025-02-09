package models

type Country struct {
	Id   int
	Name string
	Code string
}

func (c Country) GetId() int {
	return c.Id
}
