package models

type Country struct {
	Id   int
	Name string
}

func (c Country) GetId() int {
	return c.Id
}
