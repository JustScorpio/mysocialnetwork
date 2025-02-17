package repository

import (
	"database/sql"
	"user_service/internal/models"
)

type CountryRepository struct {
	db *sql.DB
}

func NewCountryRepository(db *sql.DB) *CountryRepository {
	return &CountryRepository{db: db}
}

func (r *CountryRepository) GetAll() ([]models.Country, error) {
	rows, err := r.db.Query(`
	SELECT id, name
	FROM countries
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var countries []models.Country
	for rows.Next() {
		var country models.Country
		err := rows.Scan(&country.Id, &country.Name)
		if err != nil {
			return nil, err
		}
		countries = append(countries, country)
	}

	return countries, nil
}

func (r *CountryRepository) Get(id int) (*models.Country, error) {
	var country models.Country
	err := r.db.QueryRow(`
		SELECT id, name
		FROM countries 
		WHERE id = $1
		`, id).Scan(&country.Id, &country.Name)
	if err != nil {
		return nil, err
	}
	return &country, nil
}

func (r *CountryRepository) Create(country *models.Country) error {
	err := r.db.QueryRow(`
	INSERT INTO countries (id, name) 
	VALUES ($1, $2) 
	RETURNING id
	`, country.Id, country.Name).Scan(&country.Id)
	if err != nil {
		return err
	}

	return nil
}

func (r *CountryRepository) Update(country *models.Country) error {
	_, err := r.db.Exec(`
	UPDATE countries SET name = $1
	WHERE id = $6
	`, country.Name, country.Id)
	return err
}

func (r *CountryRepository) Delete(id int) error {
	_, err := r.db.Exec(`
	DELETE FROM countries 
	WHERE id = $1
	`, id)
	return err
}
