package repository

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kk7453603/effective_mobile_golang/internal/models"
	"github.com/labstack/echo/v4"
)

var ErrUserNotFound = errors.New("user not found")
var ErrFailedToStopTask = errors.New("failed to stop task")

type SqlHandler struct {
	DB   *pgxpool.Pool
	elog echo.Logger
	dsn  string
}

func New(e echo.Logger) *SqlHandler {
	dsn := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME")
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		e.Errorf("SqlHandler init error: %v", err)
	}
	return &SqlHandler{DB: pool, dsn: dsn, elog: e}
}

func (h *SqlHandler) Migrate() {
	m, err := migrate.New(os.Getenv("DB_MIGRATIONS_PATH"), h.dsn+"?sslmode=disable")
	h.elog.Debugf("sourceURL: %s , DSN: %s", os.Getenv("DB_MIGRATIONS_PATH"), h.dsn)
	if err != nil {
		h.elog.Errorf("migration error: %v", err)
	}
	if err := m.Up(); err != nil && errors.Is(err, errors.New("migration error: no change")) {
		h.elog.Errorf("migration error: %v", err)
	}
}

func (h *SqlHandler) GetUsers(query string, args []interface{}) ([]models.User, error) {
	h.elog.Debugf("query: %s\n args:%s", query, args)
	rows, err := h.DB.Query(context.Background(), query, args...)
	if err != nil {
		h.elog.Errorf("error repository GetUsers: %s", err)
		return nil, err
	}
	defer rows.Close()
	users := []models.User{}
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.PassportNumber, &user.Surname, &user.Name, &user.Patronymic, &user.Address); err != nil {
			h.elog.Errorf("error repository GetUsers: %s", err)
			return nil, err
		}
		users = append(users, user)
		h.elog.Debugf("users: %v", users)
	}
	return users, nil
}

func (h *SqlHandler) GetUserTaskReport(passportNumber string, startDate, endDate time.Time) ([]models.TaskReport, error) {
	query := `
		SELECT t.taskname, t.content, 
		       EXTRACT(EPOCH FROM (t.endtime - t.starttime)) / 3600 AS hours, 
		       MOD(EXTRACT(EPOCH FROM (t.endtime - t.starttime)) / 60, 60) AS minutes,
		       t.starttime, t.endtime
		FROM tasks t
		JOIN users u ON t.userid = u.id
		WHERE u.passportNumber = $1 AND t.starttime >= $2 AND t.endtime <= $3
		ORDER BY EXTRACT(EPOCH FROM (t.endtime - t.starttime)) DESC
	`
	rows, err := h.DB.Query(context.Background(), query, passportNumber, startDate, endDate)
	if err != nil {
		h.elog.Errorf("error repository GetUserTaskReport: %s", err)
		return nil, err
	}
	defer rows.Close()

	var taskReports []models.TaskReport
	for rows.Next() {
		var report models.TaskReport
		if err := rows.Scan(&report.TaskName, &report.Content, &report.Hours, &report.Minutes, &report.StartTime, &report.EndTime); err != nil {
			h.elog.Errorf("error repository GetUserTaskReport: %s", err)
			return nil, err
		}
		taskReports = append(taskReports, report)
	}
	return taskReports, nil
}

func (h *SqlHandler) StartUserTask(passportNumber, taskName, content string) error {
	query := `
		INSERT INTO tasks (userid, taskname, content, starttime) 
		SELECT id, $2, $3, $4 FROM users WHERE passportNumber = $1
	`
	result, err := h.DB.Exec(context.Background(), query, passportNumber, taskName, content, time.Now())
	if err != nil {
		h.elog.Errorf("error repository StartUserTask: %s", err)
		return err
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		h.elog.Error("error repository: failed to start task")
		return errors.New("failed to start task")
	}
	return nil
}

func (h *SqlHandler) StopUserTask(passportNumber, taskName string) error {
	query := `
		UPDATE tasks
		SET endtime = $3
		FROM users
		WHERE tasks.userid = users.id
		AND users.passportNumber = $1
		AND tasks.taskname = $2
		AND tasks.endtime IS NULL
	`
	_, err := h.DB.Exec(context.Background(), query, passportNumber, taskName, time.Now())
	if err != nil {
		h.elog.Errorf("error repository StopUserTask: %s", err)
		return err
	}
	return nil
}

func (h *SqlHandler) AddUser(user models.User) error {

	var existingPassportNumber string
	row := h.DB.QueryRow(context.Background(), "SELECT passportNumber FROM users WHERE passportNumber = $1", user.PassportNumber)
	row.Scan(&existingPassportNumber)

	h.elog.Debugf("passcheck: %s", existingPassportNumber)
	if existingPassportNumber != "" {
		h.elog.Errorf("user with passport number %s already exists", user.PassportNumber)
		return errors.New("user with this passport number already exists")
	}

	query := `
		INSERT INTO users (passportNumber,surname,name,patronymic,address) 
		VALUES ($1,$2,$3,$4,$5)
	`
	_, err := h.DB.Exec(context.Background(), query, user.PassportNumber, user.Surname, user.Name, user.Patronymic, user.Address)
	if err != nil {
		h.elog.Errorf("error repository AddUser: %s", err)
		return errors.New("failed to add user")
	}
	return nil
}

func (h *SqlHandler) EditUser(passportNumber, surname, name, patronymic, address string) error {
	query := `
		UPDATE users
		SET surname = $2, name = $3, patronymic = $4, address = $5
		WHERE passportNumber = $1
	`
	result, err := h.DB.Exec(context.Background(), query, passportNumber, surname, name, patronymic, address)
	if err != nil {
		h.elog.Errorf("error repository EditUser: %s", err)
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		h.elog.Error("error repository: user not found")
		return errors.New("user not found")
	}

	return nil
}

func (h *SqlHandler) DeleteUser(passportNumber string) error {

	checkQuery := `
		SELECT COUNT(*) FROM users
		WHERE passportNumber = $1
	`
	var count int
	err := h.DB.QueryRow(context.Background(), checkQuery, passportNumber).Scan(&count)
	if err != nil {
		h.elog.Errorf("error repository DeleteUser check existence: %s", err)
		return err
	}

	if count == 0 {
		h.elog.Error("error repository: user not found")
		return ErrUserNotFound
	}

	query := `
		DELETE FROM users	
		WHERE passportNumber = $1
	`
	_, err = h.DB.Exec(context.Background(), query, passportNumber)
	if err != nil {
		h.elog.Errorf("error repository DeleteUser: %s", err)
		return err
	}

	h.elog.Debugf("repository DeleteUser: user with passport number %s deleted", passportNumber)
	return nil
}
