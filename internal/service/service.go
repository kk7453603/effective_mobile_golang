package service

import (
	"fmt"
	"time"

	"github.com/kk7453603/effective_mobile_golang/internal/models"
	"github.com/labstack/echo/v4"
)

type Repository interface {
	GetUsers(query string, args []interface{}) ([]models.User, error)
	GetUserTaskReport(passportNumber string, startDate time.Time, endDate time.Time) ([]models.TaskReport, error)
	StartUserTask(passportNumber, taskName, content string) error
	StopUserTask(passportNumber, taskName string) error
	AddUser(passportNumber, surname, name, patronymic, address string) error
	EditUser(passportNumber, surname, name, patronymic, address string) error
	DeleteUser(passportNumber string) error
}

type Service struct {
	repo Repository
	elog echo.Logger
}

func (s *Service) AddUser(passportNumber, surname, name, patronymic, address string) error {
	return s.repo.AddUser(passportNumber, surname, name, patronymic, address)
}

func (s *Service) EditUser(passportNumber, surname, name, patronymic, address string) error {
	return s.repo.EditUser(passportNumber, surname, name, patronymic, address)
}

func (s *Service) GetUserTaskStatus(passportNumber string, startDate time.Time, endDate time.Time) ([]models.TaskReport, error) {
	taskReports, err := s.repo.GetUserTaskReport(passportNumber, startDate, endDate)
	if err != nil {
		s.elog.Errorf("service GetUserTaskReport error: %s", err)
		return nil, err
	}
	return taskReports, nil
}

func (s *Service) RemoveUser(passportNumber string) error {
	return s.repo.DeleteUser(passportNumber)
}

func (s *Service) StartUserTimer(passportNumber string, taskName, content string) error {
	return s.repo.StartUserTask(passportNumber, taskName, content)
}

func (s *Service) StopUserTimer(passportNumber, taskName string) error {
	return s.repo.StopUserTask(passportNumber, taskName)
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
