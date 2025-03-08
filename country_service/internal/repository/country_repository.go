package repository

import (
	"country_service/internal/models"
	"database/sql"
)

type CountryRepository struct {
	db *sql.DB
}

func NewCountryRepository(db *sql.DB) *CountryRepository {
	return &CountryRepository{db: db}
}

func (r *CountryRepository) GetAll() ([]models.Country, error) {
	rows, err := r.db.Query("SELECT id, name, code, population FROM countries")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var countries []models.Country
	for rows.Next() {
		var country models.Country
		err := rows.Scan(&country.Id, &country.Name, &country.Code, &country.Population)
		if err != nil {
			return nil, err
		}
		countries = append(countries, country)
	}

	return countries, nil
}

func (r *CountryRepository) Get(id int) (*models.Country, error) {
	var country models.Country
	err := r.db.QueryRow("SELECT id, name, code, population FROM countries WHERE id = $1", id).Scan(&country.Id, &country.Name, &country.Code, &country.Population)
	if err != nil {
		return nil, err
	}
	return &country, nil
}

func (r *CountryRepository) Create(country *models.Country) error {
	err := r.db.QueryRow("INSERT INTO countries (name, code, population) VALUES ($1, $2) RETURNING id", country.Name, country.Code, country.Population).Scan(&country.Id)
	if err != nil {
		return err
	}

	return nil
}

func (r *CountryRepository) Update(country *models.Country) error {
	_, err := r.db.Exec("UPDATE countries SET name = $1, code = $2, population = $3 WHERE id = $4", country.Name, country.Code, country.Population, country.Id)
	return err
}

func (r *CountryRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM countries WHERE id = $1", id)
	return err
}
