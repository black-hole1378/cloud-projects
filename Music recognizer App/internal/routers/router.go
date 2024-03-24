package routers

import (
	"awesomeProject/internal/handlers"
	"github.com/labstack/echo/v4"
)

func SetupRouter(ec *echo.Echo) {
	ec.GET("/", handlers.HomeHandler)
	ec.POST("/upload", handlers.UploadHandler)
}
