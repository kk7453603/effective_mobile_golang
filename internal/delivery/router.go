package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/kk7453603/effective_mobile_golang/internal/models"
	"github.com/labstack/echo/v4"
)

type Service interface {
	GetAllCars(car models.Car, limit int) ([]models.Car, error)
	RemoveCarById(RegNum string) error
	UpdateCar(car models.Car) error
	AddCar(car models.Car) error
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
			d.logger.Debug(err)
			return c.JSONPretty(http.StatusBadRequest, &models.Response_Error{Message: "Bad request"}, "  ")
		}
		d.logger.Info(car) // временно
		if car.RegNum == "" {
			return c.JSONPretty(http.StatusNotFound, &models.Response_Error{Message: "No RegNum field in query params."}, "  ")
		}
		limit := c.QueryParams().Get("limit")
		parselim, err := strconv.Atoi(limit)
		if err != nil {
			parselim = 10
			d.logger.Debugf("limit field is empty or fail, set default value %d, error: %v", parselim, err)
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
		car := new(models.Car)
		err := c.Bind(&car)
		if err != nil {
			d.logger.Debug(err)
			return c.JSONPretty(http.StatusBadRequest, &models.Response_Error{Message: "Bad request"}, "  ")
		}
		resp_car := &models.Car{}
		if car.RegNum != "" {
			resp_car.RegNum = car.RegNum
			if car.Model != "" {
				resp_car.Model = car.Model
			}
			if resp_car.Owner.Name != "" || resp_car.Owner.Surname != "" {
				resp_car.Owner = car.Owner
			}
			if car.Mark != "" {
				resp_car.Mark = car.Mark
			}
			if car.Year != 0 {
				resp_car.Year = car.Year
			}
			err := d.serv.UpdateCar(*resp_car)
			if err != nil {
				d.logger.Debug(err)
				return c.JSONPretty(500, &models.Response_Error{Message: err.Error()}, "  ")
			}
			return c.JSONPretty(http.StatusOK, &models.Response_OK{Status: "Ok"}, "  ")
		}
		return c.JSONPretty(http.StatusNotFound, &models.Response_Error{Message: "field regnum is empty"}, "  ")
	})

	g.POST("/add_cars", func(c echo.Context) error {
		req_car := new(models.Request_Add_Cars)
		err := c.Bind(&req_car)
		if err != nil {
			d.logger.Debug(err)
			return c.JSONPretty(http.StatusBadRequest, &models.Response_Error{Message: "Bad request"}, "  ")
		}
		for _, regNum := range req_car.RegNums {
			apiURL := fmt.Sprintf("%s?regNum=%s", os.Getenv("API_URL"), regNum)
			resp, err := http.Get(apiURL)
			if err != nil {
				d.logger.Debugf("Ошибка GET запроса к API: %v", err)
				return c.JSONPretty(500, &models.Response_Error{Message: err.Error()}, "  ")
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				d.logger.Debugf("Ошибка ответа от внешнего API: %v", resp.StatusCode)
				return c.JSONPretty(500, &models.Response_Error{Message: "Ошибка ответа от внешнего API"}, "  ")
			}
			var carInfo models.Car
			if err := json.NewDecoder(resp.Body).Decode(&carInfo); err != nil {
				d.logger.Debugf("Ошибка декодирования JSON от внешенего API: %v", err)
				return c.JSONPretty(500, &models.Response_Error{Message: err.Error()}, "  ")
			}
			err = d.serv.AddCar(carInfo)
			if err != nil {
				d.logger.Debug(err)
				return c.JSONPretty(500, &models.Response_Error{Message: err.Error()}, "  ")
			}
		}
		return c.JSONPretty(http.StatusOK, &models.Response_OK{Status: "Ok"}, "  ")
	})
}
