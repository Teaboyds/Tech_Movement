package routes

import (
	"backend-tech-movement/controllers"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	api := e.Group("/api")

	api.GET("/users", controllers.GetUsers)
	api.POST("/users", controllers.CreateUser)
}
