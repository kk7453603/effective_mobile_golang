package delivery

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kk7453603/effective_mobile_golang/internal/models"
	"github.com/labstack/echo/v4"
)

var ErrUserNotFound = errors.New("user not found")

type Service interface {
	GetAllUsers(filter map[string]string, limit int, page int) ([]models.User, error)
	GetUserTaskStatus(passportNumber string, startDate time.Time, endDate time.Time) ([]models.TaskReport, error)
	StartUserTimer(passportNumber string, taskName, content string) error
	StopUserTimer(passportNumber, taskName string) error
	RemoveUser(passportNumber string) error
	EditUser(passportNumber, surname, name, patronymic, address string) error
	AddUser(user models.User) error
}

type Delivery struct {
	serv   Service
	logger echo.Logger
}

func New(s Service, lg echo.Logger) *Delivery {
	return &Delivery{serv: s, logger: lg}
}

func Validate_passport(passport string) error {
	parsed_pasport := strings.Split(passport, " ")

	if len(parsed_pasport) != 2 {
		return errors.New("please, type passport inf format: xxxx xxxxxx")
	}

	passport_serial := parsed_pasport[0]
	passport_number := parsed_pasport[1]

	if passport == "" || passport_serial == "" || passport_number == "" || len(passport_serial) != 4 || len(passport_number) != 6 {
		return errors.New("passport value is incorrect")
	}
	return nil
}

func (d *Delivery) InitRoutes(g *echo.Group) {

	g.GET("/get_users", d.GetUsers)

	g.GET("/get_user_tasks", d.GetUserTasks)

	g.POST("/start_user_task", d.StartUserTask)

	g.POST("/stop_user_task", d.StopUserTask)

	g.POST("/add_user", d.AddUser)

	g.POST("/edit_user", d.EditUser)

	g.POST("/delete_user", d.DeleteUser)
}

// @Summary		Get all users
// @Description	Get all users with optional filters
// @Accept			json
// @Produce		json
// @Param			passportNumber	query		string	false	"Passport Number"
// @Param			surname			query		string	false	"Surname"
// @Param			name			query		string	false	"Name"
// @Param			patronymic		query		string	false	"Patronymic"
// @Param			address			query		string	false	"Address"
// @Param			page			query		int		false	"Page"
// @Param			limit			query		int		false	"Limit"
// @Success		200				{array}		models.User
// @Failure		400				{object}	models.Response_Error
// @Failure		500				{object}	models.Response_Error
// @Router			/get_users [get]
func (d *Delivery) GetUsers(c echo.Context) error {
	filter := map[string]string{
		"passportNumber": c.QueryParam("passportNumber"),
		"surname":        c.QueryParam("surname"),
		"name":           c.QueryParam("name"),
		"patronymic":     c.QueryParam("patronymic"),
		"address":        c.QueryParam("address"),
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response_Error{Error: "Failed to get page value"})
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response_Error{Error: "Failed to get limit value"})
	}

	users, err := d.serv.GetAllUsers(filter, limit, page)
	if err != nil {
		d.logger.Errorf("delivery /get_users error:%s", err)
		return c.JSON(http.StatusInternalServerError, models.Response_Error{Error: "Failed to get users"})
	}
	d.logger.Debugf("users: %v", users)
	return c.JSONPretty(http.StatusOK, users, "  ")
}

// @Summary		Get user tasks
// @Description	Get user tasks for a specified period
// @Accept			json
// @Produce		json
// @Param			passportNumber	query		string	true	"Passport Number"
// @Param			start_date		query		string	true	"Start Date"	Format(date)
// @Param			end_date		query		string	true	"End Date"		Format(date)
// @Success		200				{array}		models.TaskReport
// @Failure		400				{object}	models.Response_Error
// @Failure		500				{object}	models.Response_Error
// @Router			/get_user_tasks [get]
func (d *Delivery) GetUserTasks(c echo.Context) error {
	passportNumber := c.QueryParam("passportNumber")

	if err := Validate_passport(passportNumber); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response_Error{Error: err.Error()})
	}

	startDateStr := c.QueryParam("start_date")
	endDateStr := c.QueryParam("end_date")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response_Error{Error: "Invalid start_date parameter"})
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response_Error{Error: "Invalid end_date parameter"})
	}

	taskReports, err := d.serv.GetUserTaskStatus(passportNumber, startDate, endDate)
	if err != nil {
		d.logger.Errorf("delivery /get_user_task_report error: %s", err)
		return c.JSON(http.StatusInternalServerError, models.Response_Error{Error: "Failed to get user task report"})
	}

	d.logger.Debugf("tasks: %v", taskReports)
	return c.JSONPretty(http.StatusOK, taskReports, "  ")
}

// @Summary		Start user task
// @Description	Start a new task for a user
// @Accept			application/x-www-form-urlencoded
// @Produce		json
// @Param			passportNumber	formData	string	true	"Passport Number"
// @Param			task_name		formData	string	true	"Task Name"
// @Param			content			formData	string	true	"Content"
// @Success		200				{object}	models.Response_OK
// @Failure		400				{object}	models.Response_Error
// @Failure		500				{object}	models.Response_Error
// @Router			/start_user_task [post]
func (d *Delivery) StartUserTask(c echo.Context) error {
	passportNumber := c.FormValue("passportNumber")

	if err := Validate_passport(passportNumber); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response_Error{Error: err.Error()})
	}

	taskName := c.FormValue("task_name")
	content := c.FormValue("content")

	if taskName == "" || content == "" {
		return c.JSON(http.StatusBadRequest, models.Response_Error{Error: "task_name, and content are required"})
	}

	err := d.serv.StartUserTimer(passportNumber, taskName, content)
	if err != nil {
		d.logger.Errorf("delivery /start_user_task error: %s", err)
		return c.JSON(http.StatusInternalServerError, models.Response_Error{Error: "Failed to start user task"})
	}

	return c.JSON(http.StatusOK, models.Response_OK{Status: "Task started"})
}

// @Summary		Stop user task
// @Description	Stop an ongoing task for a user
// @Accept			application/x-www-form-urlencoded
// @Produce		json
// @Param			passportNumber	formData	string	true	"Passport Number"
// @Param			task_name		formData	string	true	"Task Name"
// @Success		200				{object}	models.Response_OK
// @Failure		400				{object}	models.Response_Error
// @Failure		500				{object}	models.Response_Error
// @Router			/stop_user_task [post]
func (d *Delivery) StopUserTask(c echo.Context) error {
	passportNumber := c.FormValue("passportNumber")
	taskName := c.FormValue("task_name")

	if passportNumber == "" || taskName == "" {
		return c.JSON(http.StatusBadRequest, models.Response_Error{Error: "passportNumber and task_name are required"})
	}

	if err := Validate_passport(passportNumber); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response_Error{Error: err.Error()})
	}

	err := d.serv.StopUserTimer(passportNumber, taskName)
	if err != nil {
		d.logger.Errorf("delivery /stop_user_task error: %s", err)
		return c.JSON(http.StatusInternalServerError, models.Response_Error{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, models.Response_OK{Status: "Task stopped"})
}

// @Summary		Add new user
// @Description	Add a new user by passport number
// @Accept			application/x-www-form-urlencoded
// @Produce		json
// @Param			passportNumber	formData	string	true	"Passport Number"
// @Success		200				{object}	models.Response_OK
// @Failure		400				{object}	models.Response_Error
// @Failure		500				{object}	models.Response_Error
// @Router			/add_user [post]
func (d *Delivery) AddUser(c echo.Context) error {
	passport := c.FormValue("passportNumber")

	if err := Validate_passport(passport); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response_Error{Error: err.Error()})
	}

	parsed_pasport := strings.Split(passport, " ")
	passport_serial := parsed_pasport[0]
	passport_number := parsed_pasport[1]

	query := fmt.Sprintf("%s?passportSerie=%s&passportNumber=%s", os.Getenv("API_URL"), passport_serial, passport_number)
	d.logger.Debugf("API query prepare: %s", query)
	res, err := http.Get(query)
	if err != nil {
		d.logger.Errorf("API request error: %s", err)
		return c.JSON(http.StatusInternalServerError, models.Response_Error{Error: err.Error()})
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		d.logger.Errorf("API request body parse error: %s", err)
		return c.JSON(http.StatusInternalServerError, models.Response_Error{Error: err.Error()})
	}
	var user models.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		d.logger.Errorf("API request json unmarshal error: %s", err)
		return c.JSON(http.StatusInternalServerError, models.Response_Error{Error: err.Error()})
	}
	user.PassportNumber = passport

	d.logger.Debugf("user unmarshal: %v", user)

	err = d.serv.AddUser(user)
	if err != nil {
		d.logger.Errorf("delivery /add_user error: %s", err)
		return c.JSON(http.StatusInternalServerError, models.Response_Error{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, models.Response_OK{Status: "ok"})
}

// @Summary		Edit user details
// @Description	Edit user details by passport number
// @Accept			application/x-www-form-urlencoded
// @Produce		json
// @Param			passportNumber	formData	string	true	"Passport Number"
// @Param			surname			formData	string	false	"Surname"
// @Param			name			formData	string	false	"Name"
// @Param			patronymic		formData	string	false	"Patronymic"
// @Param			address			formData	string	false	"Address"
// @Success		200				{object}	models.Response_OK
// @Failure		400				{object}	models.Response_Error
// @Failure		500				{object}	models.Response_Error
// @Router			/edit_user [post]
func (d *Delivery) EditUser(c echo.Context) error {
	user := models.Response_User{
		PassportNumber: c.FormValue("passportNumber"),
		Surname:        c.FormValue("surname"),
		Name:           c.FormValue("name"),
		Patronymic:     c.FormValue("patronymic"),
		Address:        c.FormValue("address"),
	}

	if err := Validate_passport(user.PassportNumber); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response_Error{Error: err.Error()})
	}

	d.logger.Debugf("user data: %v", user)
	err := d.serv.EditUser(user.PassportNumber, user.Surname, user.Name, user.Patronymic, user.Address)
	if err != nil {
		d.logger.Errorf("delivery /edit_user error: %s", err)
		return c.JSON(http.StatusInternalServerError, models.Response_Error{Error: "Failed to edit user"})
	}

	return c.JSON(http.StatusOK, models.Response_OK{Status: "User updated"})
}

// @Summary		Delete user
// @Description	Delete user by passport number
// @Accept			application/x-www-form-urlencoded
// @Produce		json
// @Param			passportNumber	formData	string	true	"Passport Number"
// @Success		200				{object}	models.Response_OK
// @Failure		400				{object}	models.Response_Error
// @Failure		500				{object}	models.Response_Error
// @Router			/delete_user [post]
func (d *Delivery) DeleteUser(c echo.Context) error {
	passportNumber := c.FormValue("passportNumber")

	if err := Validate_passport(passportNumber); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response_Error{Error: err.Error()})
	}

	err := d.serv.RemoveUser(passportNumber)
	d.logger.Debugf("DeleteUser delivery level err:%v", err)
	if errors.Is(err, ErrUserNotFound) {
		return c.JSON(http.StatusBadRequest, models.Response_Error{Error: err.Error()})
	} else if err != nil {
		d.logger.Errorf("delivery /delete_user error: %s", err)
		return c.JSON(http.StatusInternalServerError, models.Response_Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, models.Response_OK{Status: "User deleted"})
}
