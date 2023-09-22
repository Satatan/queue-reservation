package main

import (
	"fmt"

	"github.com/labstack/echo/v4"

	"queue_reservation/app/api/factory"
	"queue_reservation/configs"

	_ "github.com/swaggo/echo-swagger/example/docs"
)

// @title Swagger QR API
// @version 1.0
// @description Open API documentation for Queue Reservation
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
func main() {

	var (
		conf = configs.Init()
	)

	e := echo.New()

	factory.QRDependencyResolve(e)

	// Start the server
	e.Start(fmt.Sprintf(":%s", conf.App.Port))
}
