package services

import (
	"mysocialnetwork/user_service/internal/models"
	"mysocialnetwork/user_service/internal/repository"
)

type CrudService[T models.Entity] struct {
	repo repository.Repository[T]
}

func NewCrudService[T models.Entity](repo repository.Repository[T]) *CrudService[T] {
	return &CrudService[T]{repo: repo}
}

func (s *CrudService[T]) GetAll() ([]T, error) {
	return s.repo.GetAll()
}

func (s *CrudService[T]) Get(id int) (*T, error) {
	return s.repo.Get(id)
}

func (s *CrudService[T]) Create(entity *T) error {
	return s.repo.Create(entity)
}

func (s *CrudService[T]) Update(entity *T) error {
	return s.repo.Update(entity)
}

func (s *CrudService[T]) Delete(id int) error {
	return s.repo.Delete(id)
}
