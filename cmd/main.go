package main

import (
	"github.com/kk7453603/effective_mobile_golang/internal/repository"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	sql_handler := repository.New(e)
	err := sql_handler.Migrate()
	if err != nil {
		e.Logger.Fatal(err)
	}

}
