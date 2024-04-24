package repository

import (
	"context"
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
	elog *echo.Echo //немного не "чистая архитектура"
	dsn  string
}

func New(e *echo.Echo) *SqlHandler {
	err := godotenv.Load()
	if err != nil {
		e.Logger.Errorf(".env load error: %v", err)
	}
	dsn := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME")
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		e.Logger.Errorf("SqlHandler init error: %v", err)
	}
	return &SqlHandler{DB: pool, dsn: dsn, elog: e}
}

func (h *SqlHandler) Migrate() {
	m, err := migrate.New(os.Getenv("DB_MIGRATIONS_PATH"), h.dsn+"?sslmode=disable")
	if err != nil {
		h.elog.Logger.Errorf("migration error: %v", err)
	}
	if err := m.Up(); err != nil {
		h.elog.Logger.Errorf("migration error: %v", err)
	}
}

func (h *SqlHandler) GetAllCars() []models.Car {
	var cars []models.Car = make([]models.Car, 0, 10)
	car := &models.Car{}

	rows, err := h.DB.Query(context.Background(), "SELECT * FROM cars;")
	if err != nil {
		h.elog.Logger.Fatal(err)
	}
	for rows.Next() {

		err := rows.Scan(&car)
		if err != nil {
			h.elog.Logger.Fatal(err)
		}
		cars = append(cars, *car)
	}
	return cars
}
