package service

import (
	"github.com/kk7453603/effective_mobile_golang/internal/models"
)

type Repository interface {
	GetAllCars(prep_query string, filterFields []string, limit int) ([]models.User, error)
}

type Service struct {
	repo Repository
}

func New(rep Repository) *Service {
	return &Service{repo: rep}
}

func (s *Service) GetAllUsers(filter map[string]string, limit int) ([]models.User, error) {

}
