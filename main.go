package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)
func main() {
	e := echo.New()

	e.Use(middleware.Logger())

	e.GET("/", Default)
	e.POST("/check_response", CheckResponse)
	e.POST("/return_response", ReturnResponse)

	e.Logger.Fatal(e.Start(":3030"))
}

