package delivery

import (
	"net/http"
	"strconv"
	"time"

	"github.com/kk7453603/effective_mobile_golang/internal/models"
	"github.com/labstack/echo/v4"
)

type Service interface {
	GetAllUsers(filter map[string]string, limit int, page int) ([]models.User, error)
	GetUserTaskStatus(passportNumber string, startDate time.Time, endDate time.Time) ([]models.TaskReport, error)
	StartUserTimer(passportNumber string) error
	StopUserTimer(passportNumber string) error
	RemoveUser(passportNumber string) error
	EditUser(passportNumber string) error
	AddUser(passportNumber string) error
}

type Delivery struct {
	serv   Service
	logger echo.Logger
}

func New(s Service, lg echo.Logger) *Delivery {
	return &Delivery{serv: s, logger: lg}
}

func (d *Delivery) InitRoutes(g *echo.Group) {

	g.GET("/get_users", func(c echo.Context) error {
		filter := map[string]string{
			"passportNumber": c.QueryParam("passportNumber"),
			"surname":        c.QueryParam("surname"),
			"name":           c.QueryParam("name"),
			"patronymic":     c.QueryParam("patronymic"),
			"address":        c.QueryParam("address"),
		}
		page, err := strconv.Atoi(c.QueryParam("page"))
		if err != nil {
			c.JSON(http.StatusBadRequest, "page value error: "+err.Error())
		}
		limit, err := strconv.Atoi(c.QueryParam("limit"))
		if err != nil {
			c.JSON(http.StatusBadRequest, "limit value error: "+err.Error())
		}

		users, err := d.serv.GetAllUsers(filter, limit, page)
		if err != nil {
			d.logger.Errorf("delivery /get_users error:%s", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get users"})
		}
		d.logger.Debugf("users: %v", users)
		return c.JSONPretty(http.StatusOK, users, "  ")
	})

	g.GET("/get_user", func(c echo.Context) error {
		passportNumber := c.QueryParam("passport_number")
		if passportNumber == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "passport_number parameter is required"})
		}

		startDateStr := c.QueryParam("start_date")
		endDateStr := c.QueryParam("end_date")

		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid start_date parameter"})
		}

		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid end_date parameter"})
		}

		taskReports, err := d.serv.GetUserTaskStatus(passportNumber, startDate, endDate)
		if err != nil {
			d.logger.Errorf("delivery /get_user_task_report error: %s", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get user task report"})
		}

		d.logger.Debugf("tasks: %v", taskReports)
		return c.JSONPretty(http.StatusOK, taskReports, "  ")
	})
}
