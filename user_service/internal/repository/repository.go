package repository

import "user-service/internal/models"

type Repository[T models.Entity] interface {
	GetAll() ([]T, error)
	Get(id int) (*T, error)
	Create(entity *T) error
	Update(entity *T) error
	Delete(id int) error
}
