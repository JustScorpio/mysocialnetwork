package services

import (
	"country_service/internal/models"
	"country_service/internal/repository"
)

type Service[T models.Entity] struct {
	repo repository.Repository[T]
}

func NewService[T models.Entity](repo repository.Repository[T]) *Service[T] {
	return &Service[T]{repo: repo}
}

func (s *Service[T]) GetAll() ([]T, error) {
	return s.repo.GetAll()
}

func (s *Service[T]) Get(id int) (*T, error) {
	return s.repo.Get(id)
}

func (s *Service[T]) Create(entity *T) error {
	return s.repo.Create(entity)
}

func (s *Service[T]) Update(entity *T) error {
	return s.repo.Update(entity)
}

func (s *Service[T]) Delete(id int) error {
	return s.repo.Delete(id)
}
