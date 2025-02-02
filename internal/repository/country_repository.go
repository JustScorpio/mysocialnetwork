package repository

import (
	"database/sql"
	"network/internal/models"
)

type CountryRepository struct {
	db *sql.DB
}

func NewCountryRepository(db *sql.DB) Repository[models.Country] {
	return &CountryRepository{db: db}
}

func (r *CountryRepository) GetAll() ([]models.Country, error) {
	rows, err := r.db.Query("SELECT id, name FROM countries")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var countries []models.Country
	for rows.Next() {
		var country models.Country
		if err := rows.Scan(country.Id, country.Name); err != nil {
			return nil, err
		}
		countries = append(countries, country)
	}

	return countries, nil
}

func (r *CountryRepository) Get(id int) (*models.Country, error) {
	var country models.Country
	err := r.db.QueryRow("SELECT id, name FROM countries WHERE id = ?", id).Scan(&country.Id, &country.Name)
	if err != nil {
		return nil, err
	}
	return &country, nil
}

func (r *CountryRepository) Create(country *models.Country) error {
	result, err := r.db.Exec("INSERT INTO countries (name) VALUES (?)", country.Name)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	country.Id = int(id)
	return nil
}

func (r *CountryRepository) Update(country *models.Country) error {
	_, err := r.db.Exec("UPDATE countries SET name = ?, WHERE id = ?", country.Name, country.Id)
	return err
}

func (r *CountryRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM countries WHERE id = ?", id)
	return err
}
