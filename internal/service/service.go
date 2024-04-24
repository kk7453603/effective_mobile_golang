package service

import "github.com/kk7453603/effective_mobile_golang/internal/models"

type Repository interface {
	GetAllCars() []models.Car
	GetCarByRegNum(regnum string) string
	RemoveCarById(carid uint64) string
	UpdateCar(RegNum string, Mark string, Model string, Year string) string
	AddCar(RegNums []string) string
}

type Service struct {
	repo Repository
}
