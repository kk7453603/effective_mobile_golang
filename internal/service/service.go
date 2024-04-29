package service

import (
	"github.com/kk7453603/effective_mobile_golang/internal/models"
)

type Repository interface {
	GetAllCars(limit int) []models.Car
	GetCarByRegNum(regnum string) models.Car
	RemoveCarById(RegNum string) error
	UpdateCar(RegNum string, Mark string, Model string, Owner uint64, Year string) error
	AddCar(RegNums []string) string
}

type Service struct {
	repo Repository
}

func New(rep Repository) *Service {
	return &Service{repo: rep}
}

/*
func (s *Service) GetFilteredCars(fields []string,limit int) []models.Response_Car {

	cars:=s.repo.GetAllCars(limit)

	resp_car:=models.Response_Car{}
	checkfields:=[]string{""}
	for _,field:= range fields {
		switch strings.ToLower(field) {
		case "regnum":
			b
		}
	}
}
*/

func (s *Service) RemoveCarById(RegNum string) error {
	return s.repo.RemoveCarById(RegNum)
}

func (s *Service) GetAllCars(fields []string, limit int) []models.Response_Car {

}
