package repository

import (
	"context"
	"database/sql"
	"weather-client/internal/pkg/models"
)

type CityRepository interface {
	GetAllCities(ctx context.Context) ([]models.City, error)
}

type cityRepository struct {
	db *sql.DB
}

func NewCityRepository(db *sql.DB) CityRepository {
	return &cityRepository{db: db}
}

func (r *cityRepository) GetAllCities(ctx context.Context) ([]models.City, error) {

	rows, err := r.db.QueryContext(ctx,
		"SELECT euid, city_name, country, latitude, longitude FROM cities")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cities []models.City

	for rows.Next() {
		var c models.City
		err := rows.Scan(
			&c.EUID,
			&c.CityName,
			&c.Country,
			&c.Latitude,
			&c.Longitude,
		)
		if err != nil {
			return nil, err
		}
		cities = append(cities, c)
	}

	return cities, nil
}
