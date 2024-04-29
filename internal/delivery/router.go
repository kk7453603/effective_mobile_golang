package delivery

import (
	"net/http"
	"strconv"

	"github.com/kk7453603/effective_mobile_golang/internal/models"
	"github.com/labstack/echo/v4"
)

type Service interface {
	GetAllCars(car models.Car, limit int) ([]models.Car, error)
	GetCarByRegNum(regnum string) models.Car
	RemoveCarById(RegNum string) error
	UpdateCar(RegNum string, Mark string, Model string, Owner uint64, Year string) error
	AddCar(RegNums []string) string
}

type Delivery struct {
	serv   Service
	logger echo.Logger //убрать логгер с других слоев
}

func New(s Service, lg echo.Logger) *Delivery {
	return &Delivery{serv: s, logger: lg}
}

func (d *Delivery) InitRoutes(g *echo.Group) {
	g.GET("/get_cars", func(c echo.Context) error {
		car := new(models.Car)
		err := c.Bind(&car)
		if err != nil {
			return c.String(http.StatusBadRequest, "Bad request")
		}
		d.logger.Info(car) // временно
		if car.RegNum == "" {
			return c.String(http.StatusNotFound, "No RegNum field in query params.")
		}
		limit := c.QueryParams().Get("limit")
		parselim, err := strconv.Atoi(limit)
		if err != nil {
			d.logger.Debugf("limit field is empty or fail: %v", err)
			parselim = 10
		}

		resp, err := d.serv.GetAllCars(*car, parselim)
		if err != nil {
			d.logger.Debugf("GetAllCars error: %v", err)
			return c.JSONPretty(500, &models.Response_Error{Message: err.Error()}, "  ")
		}
		return c.JSONPretty(http.StatusOK, resp, "  ")
	})

	g.DELETE("/remove_car", func(c echo.Context) error {
		regnum := c.FormValue("regnum")
		if regnum != "" {
			err := d.serv.RemoveCarById(regnum)
			if err != nil {
				d.logger.Debugf("Error on %s in RemoveCarById: %v", c.Request().URL.String(), err)
				return c.JSONPretty(500, &models.Response_Error{Message: err.Error()}, "  ")
			}
			return c.JSONPretty(http.StatusOK, &models.Response_OK{Status: "Ok"}, "  ")
		}
		return c.JSONPretty(http.StatusNotFound, &models.Response_Error{Message: "field regnum is empty"}, "  ")
	})

	g.PUT("/update_car", func(c echo.Context) error {
		regnum := c.FormValue("regnum")
		if regnum != "" {

		}
	})
}
