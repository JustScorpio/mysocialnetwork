package repository

import (
	"database/sql"
	"user_service/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) Repository[models.User] {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetAll() ([]models.User, error) {
	rows, err := r.db.Query("SELECT id, username, mail FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Id, &user.UserName, &user.Mail)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) Get(id int) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow("SELECT id, username, mail FROM users WHERE id = $1", id).Scan(&user.Id, &user.UserName, &user.Mail)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(user *models.User) error {
	err := r.db.QueryRow("INSERT INTO users (username, mail) VALUES ($1, $2) RETURNING id", user.UserName, user.Mail).Scan(&user.Id)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Update(user *models.User) error {
	_, err := r.db.Exec("UPDATE users SET username = $1, mail = $2 WHERE id = $3", user.UserName, user.Mail, user.Id)
	return err
}

func (r *UserRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}
