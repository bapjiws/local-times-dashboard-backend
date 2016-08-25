package datastore

import "timezones_mc/revel_app/app/models"

//TODO: Create type Storage that will embed all the particular storages?
type CityStorage interface {
	AddCity(city *models.City) error
}
