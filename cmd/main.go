package main

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/kk7453603/effective_mobile_golang/docs"
	"github.com/kk7453603/effective_mobile_golang/internal/delivery"
	"github.com/kk7453603/effective_mobile_golang/internal/repository"
	"github.com/kk7453603/effective_mobile_golang/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"
)

//	@title			Time Tracker API
//	@version		1.0
//	@description	This is a time tracker server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		myapp.local
// @BasePath	/
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("godotenv error: %v", err)
	}
	e := echo.New()
	e.Logger.Info("Переменные среды загружены")
	if os.Getenv("DEBUG") == "on" {
		e.Debug = true
		e.Logger.SetLevel(log.DEBUG)
		e.Logger.Info("DEBUG режим включен")
	}

	sql_handler := repository.New(e.Logger)
	sql_handler.Migrate()

	e.Use(middleware.Logger())
	//e.Use(middleware.Recover())

	// Swagger documentation route
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Мок-обработчик для эндпоинта /info
	e.GET("/info", func(c echo.Context) error {
		passportSerie := c.QueryParam("passportSerie")
		passportNumber := c.QueryParam("passportNumber")

		// Пример данных о человеке
		people := map[string]interface{}{
			"surname":    "Иванов",
			"name":       "Иван",
			"patronymic": "Иванович",
			"address":    "г. Москва, ул. Ленина, д. 5, кв. 1",
		}

		// Проверка наличия необходимых параметров
		if passportSerie == "" || passportNumber == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Missing required parameters"})
		}

		// Возвращаем данные о человеке
		return c.JSON(http.StatusOK, people)
	})

	g := e.Group("")
	serv := service.New(sql_handler, e.Logger)
	deliv := delivery.New(serv, e.Logger)
	deliv.InitRoutes(g)

	e.Logger.Infof("service starts on: %s", e.Server.Addr)

	if err = e.Start(os.Getenv("Service_Url")); err != nil {
		e.Logger.Fatalf("Ошибка запуска сервера: %v", err)
	}
	e.Logger.Info("service start")
}
