package repository

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"

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

func (h *SqlHandler) GetAllCars(prep_query string, filterFields []string, limit int) ([]models.Car, error) {
	var cars []models.Car = make([]models.Car, 0, 10)
	query := fmt.Sprintf("SELECT %s FROM cars LIMIT %d;", prep_query, limit)
	h.elog.Debug(query)
	rows, err := h.DB.Query(context.Background(), query)
	if err != nil {
		h.elog.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	var values []interface{}
	for rows.Next() {
		values = make([]interface{}, len(filterFields))
		for i := range filterFields {
			values[i] = new(interface{})
		}
		err := rows.Scan(values...)
		if err != nil {
			h.elog.Fatal(err)
			return nil, err
		}
		car := models.Car{}
		for i, columnName := range filterFields {
			switch columnName {
			case "regnum":
				regNum, ok := values[i].(*string)
				h.elog.Info(reflect.TypeOf(values[i]))
				if ok {
					h.elog.Info(regNum)
					car.RegNum = *regNum
				}

			case "mark":
				mark, ok := values[i].(*string)
				if ok {
					car.Mark = *mark
				}

			case "model":
				model, ok := values[i].(*string)
				if ok {
					car.Model = *model
				}

			case "year":
				year, ok := values[i].(*int)
				if ok {
					car.Year = *year
				}
			default:

			}
		}

		cars = append(cars, car)
	}
	return cars, nil
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

func (h *SqlHandler) AddCar(car models.Car) error {
	// Mock на внешнее API

	//query := "INSERT INTO cars (regnum) VALUES ($1);"
	return nil

}

func (h *SqlHandler) AddUser(user models.People) error {
	return nil
}
func (h *SqlHandler) RemoveCarById(RegNum string) error {
	_, err := h.DB.Exec(context.Background(), "DELETE FROM cars WHERE regnum=$1", RegNum)
	if err != nil {
		h.elog.Fatal(err)
		return err
	}
	return nil
}

func (h *SqlHandler) UpdateCar(car_fields string, owner_fileds string) error {
	/*
		_, err := h.DB.Exec(context.Background(), "Update cars SET mark=$1, model=$2,", RegNum) //доделать
		if err != nil {
			h.elog.Fatal(err)
			return err
		}
	*/
	return nil

}
