package repository

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/kk7453603/effective_mobile_golang/internal/models"
	"github.com/labstack/echo/v4"
)

type SqlHandler struct {
	DB   *pgxpool.Pool
	elog echo.Logger //немного не "чистая архитектура"
	dsn  string
}

func New(e echo.Logger) *SqlHandler {
	err := godotenv.Load()
	if err != nil {
		e.Errorf(".env load error: %v", err)
	}
	dsn := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME")
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		e.Errorf("SqlHandler init error: %v", err)
	}
	return &SqlHandler{DB: pool, dsn: dsn, elog: e}
}

func (h *SqlHandler) Migrate() {
	m, err := migrate.New(os.Getenv("DB_MIGRATIONS_PATH"), h.dsn+"?sslmode=disable")
	if err != nil {
		h.elog.Errorf("migration error: %v", err)
	}
	if err := m.Up(); err != nil && errors.Is(err, errors.New("migration error: no change")) {
		h.elog.Errorf("migration error: %v", err)
	}
}

func (h *SqlHandler) GetAllCars(limit int) []models.Car {
	var cars []models.Car = make([]models.Car, 0, 10)
	var query string
	if limit == 0 {
		query = "SELECT owner,regnum,mark,model,year FROM cars;"
	} else {
		query = "SELECT owner,regnum,mark,model,year FROM cars LIMIT " + fmt.Sprintf("%d", limit) + ";"
	}
	rows, err := h.DB.Query(context.Background(), query)
	if err != nil {
		h.elog.Fatal(err)
	}
	for rows.Next() {
		car := &models.Car{}
		err := rows.Scan(&car.Owner, &car.RegNum, &car.Mark, &car.Model, &car.Year)
		if err != nil {
			h.elog.Fatal(err)
		}
		cars = append(cars, *car)
	}
	return cars
}

func (h *SqlHandler) GetCarByRegNum(regnum string) models.Car {
	car := &models.Car{}
	row := h.DB.QueryRow(context.Background(), "SELECT regnum,mark,model,year FROM cars WHERE regnum=$1;", regnum)
	err := row.Scan(&car.RegNum, &car.Mark, &car.Model, &car.Year)
	if err != nil {
		h.elog.Fatal(err)
	}
	return *car
}

func (h *SqlHandler) AddCar(RegNum string) {
	// Mock на внешнее API

	//query := "INSERT INTO cars (regnum) VALUES ($1);"

}
