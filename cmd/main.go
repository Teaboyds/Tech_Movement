package main

import (
	database "backend-tech-movement/config"
	"backend-tech-movement/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {

	db := database.ConnectDB()

	app := fiber.New()

	routes.SetupNewsRoutes(app,db)
	routes.SetupCategoryRoutes(app,db)


	app.Listen(":7500")
}
