// controllers/user_controller.go
package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func GetUsers(c *fiber.Ctx) error {
	user := []string{"John" , "Jane"}
	return c.JSON(fiber.Map{
		"status":fiber.StatusCreated ,
		"data": user,
	})
}

func CreateUser(c *fiber.Ctx) error {
	user := new(struct {
		Name string `json:"name"`
	})
	if err := c.BodyParser(&user); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
		})			}

	return c.JSON(fiber.Map{
		"status":fiber.StatusCreated ,
		"data": user,
	})
}
