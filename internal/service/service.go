package service

import (
	"fmt"

	"github.com/kk7453603/effective_mobile_golang/internal/models"
	"github.com/labstack/echo/v4"
)

type Repository interface {
	GetUsers(query string, args []interface{}) ([]models.User, error)
}

type Service struct {
	repo Repository
	elog echo.Logger
}

// AddUser implements delivery.Service.
func (s *Service) AddUser(passportNumber string) error {
	return nil
}

// EditUser implements delivery.Service.
func (s *Service) EditUser(passportNumber string) error {
	return nil
}

// GetUserStatus implements delivery.Service.
func (s *Service) GetUserStatus(passportNumber string) error {
	return nil
}

// RemoveUser implements delivery.Service.
func (s *Service) RemoveUser(passportNumber string) error {
	return nil
}

// StartUserTimer implements delivery.Service.
func (s *Service) StartUserTimer(passportNumber string) error {
	return nil
}

// StopUserTimer implements delivery.Service.
func (s *Service) StopUserTimer(passportNumber string) error {
	return nil
}

func New(rep Repository, loger echo.Logger) *Service {
	return &Service{repo: rep, elog: loger}
}

func (s *Service) GetAllUsers(filter map[string]string, limit int, page int) ([]models.User, error) {
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	query := "SELECT id, passportNumber, surname, name, patronymic, address FROM users WHERE 1=1"
	args := []interface{}{}
	argID := 1
	for k, v := range filter {
		if v != "" {
			query += fmt.Sprintf(" AND %s = $%d", k, argID)
			args = append(args, v)
			argID++
		}
	}
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argID, argID+1)
	args = append(args, limit, offset)
	users, err := s.repo.GetUsers(query, args)
	s.elog.Debugf("service getallusers users:%v", users)
	if err != nil {
		return nil, err
	}

	return users, nil
}
