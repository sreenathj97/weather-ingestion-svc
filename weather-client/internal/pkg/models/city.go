package models

import "github.com/google/uuid"

type City struct {
	EUID      uuid.UUID
	CityName  string
	Country   string
	Latitude  float64
	Longitude float64
}
