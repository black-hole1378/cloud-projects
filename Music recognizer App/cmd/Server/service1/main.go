package main

import (
	"awesomeProject/internal/routers"
	"github.com/labstack/echo/v4"
)

func main() {
	port := ":4000"
	ec := echo.New()
	routers.SetupRouter(ec)
	ec.Logger.Fatal(ec.Start(port))
}
