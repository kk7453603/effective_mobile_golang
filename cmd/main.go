package main

import (
	"os"

	"github.com/kk7453603/effective_mobile_golang/internal/repository"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	sql_handler := repository.New(e.Logger)
	sql_handler.Migrate()
	e.GET("/", func(c echo.Context) error {
		test := []string{"asas", "asasas"}
		sql_handler.AddCar(test)
		return c.JSON(200, sql_handler.GetAllCars(0))
	})
	e.Start(os.Getenv("Service_Url"))
}
