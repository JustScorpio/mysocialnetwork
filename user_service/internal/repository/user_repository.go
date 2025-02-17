package repository

import (
	"database/sql"
	"fmt"
	"user_service/internal/models"
)

type UserRepository struct {
	db               *sql.DB
	countryCacheRepo *CountryRepository
}

func NewUserRepository(db *sql.DB, countryCacheRepo *CountryRepository) *UserRepository {
	return &UserRepository{
		db:               db,
		countryCacheRepo: countryCacheRepo,
	}
}

func (r *UserRepository) GetAllu() ([]models.User, error) {
	rows, err := r.db.Query(`
	SELECT id, username, mail 
	FROM users
	`)
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

func (r *UserRepository) GetAll() ([]models.User, error) {
	rows, err := r.db.Query(`
        SELECT 
            users.id, 
            users.username, 
            users.name, 
            users.passwordhash, 
            users.mail, 
            countries.id AS country_id, 
            countries.name AS country_name
        FROM users
        LEFT JOIN countries ON users.country = countries.id
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		var countryID sql.NullInt64
		var countryName sql.NullString

		err := rows.Scan(
			&user.Id,
			&user.UserName,
			&user.Name,
			&user.PasswordHash,
			&user.Mail,
			&countryID,
			&countryName,
		)
		if err != nil {
			return nil, err
		}

		if countryID.Valid {
			user.Country = &models.Country{
				Id:   int(countryID.Int64),
				Name: countryName.String,
			}
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) Get(id int) (*models.User, error) {
	var user models.User
	var countryId *int

	err := r.db.QueryRow(`
		SELECT 
			users.id, 
			users.username, 
			users.name, 
			users.passwordhash, 
			users.mail
			users.country
		FROM users
		WHERE id = $1
		`, id).Scan(
		&user.Id,
		&user.UserName,
		&user.Name,
		&user.PasswordHash,
		&user.Mail,
		&countryId,
	)
	if err != nil {
		return nil, err
	}

	// Если countryId указан, получаем данные о стране
	if countryId != nil {
		country, err := r.countryCacheRepo.Get(*countryId)
		if err != nil {
			return nil, fmt.Errorf("failed to get country: %w", err)
		}
		user.Country = country
	}

	return &user, nil
}

func (r *UserRepository) Create(user *models.User) error {
	// Проверка существования страны
	var countryId *int = nil
	if user.Country != nil {
		_, err := r.countryCacheRepo.Get(user.Country.Id)
		if err != nil {
			return fmt.Errorf("failed to get country: %w", err)
		}
		countryId = &user.Country.Id
	}

	err := r.db.QueryRow(`
	INSERT INTO users (username, name, passwordhash, mail, country) 
	VALUES ($1, $2, $3, $4, $5) 
	RETURNING id
	`, user.UserName, user.Name, user.PasswordHash, user.Mail, countryId).Scan(&user.Id)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Update(user *models.User) error {
	// Проверка существования страны
	var countryId *int = nil
	if user.Country != nil {
		_, err := r.countryCacheRepo.Get(user.Country.Id)
		if err != nil {
			return fmt.Errorf("failed to get country: %w", err)
		}
		countryId = &user.Country.Id
	}

	_, err := r.db.Exec(`
	UPDATE users SET username = $1, name = $2, passwordhash = $3, mail = $4, country = $5 
	WHERE id = $6
	`, user.UserName, user.Name, user.PasswordHash, user.Mail, countryId, user.Id)
	return err
}

func (r *UserRepository) Delete(id int) error {
	_, err := r.db.Exec(`
	DELETE FROM users 
	WHERE id = $1`, id)
	return err
}
