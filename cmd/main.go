package main

import (
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/kk7453603/effective_mobile_golang/internal/delivery"
	"github.com/kk7453603/effective_mobile_golang/internal/models"
	"github.com/kk7453603/effective_mobile_golang/internal/repository"
	"github.com/kk7453603/effective_mobile_golang/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()
	e.Debug = true
	e.Logger.SetLevel(log.DEBUG)
	e.Use(middleware.Logger())
	g := e.Group("")
	sql_handler := repository.New(e.Logger)
	sql_handler.Migrate()
	serv := service.New(sql_handler)
	deliv := delivery.New(serv, e.Logger)
	deliv.InitRoutes(g)
	/*
		e.GET("/", func(c echo.Context) error {
			test := "try"
			sql_handler.AddCar(test)
			return c.JSON(200, sql_handler.GetAllCars(0))
		})*/
	// Mock API
	e.GET("/info", func(c echo.Context) error {
		regNum := c.QueryParam("regNum")
		if regNum == "" {
			return c.JSON(http.StatusBadRequest, "Bad request")
		}

		car := models.Car{
			RegNum: regNum,
			Mark:   "TestMark",
			Model:  "TestModel",
			Year:   2022,
			Owner: models.People{
				Name:       "John",
				Surname:    "Doe",
				Patronymic: "Smith",
			},
		}

		return c.JSON(http.StatusOK, car)
	})

	e.GET("/add", func(c echo.Context) error {
		regNum := c.QueryParam("regNum")
		if regNum == "" {
			return c.JSON(http.StatusBadRequest, "Bad request")
		}
		req := httptest.NewRequest(http.MethodGet, "/info?regNum="+regNum, nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		res := rec.Body.Bytes()
		return c.JSONBlob(200, res)

	})

	e.Start(os.Getenv("Service_Url"))
}
