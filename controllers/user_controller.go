// controllers/user_controller.go
package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetUsers(c echo.Context) error {
	return c.JSON(http.StatusOK, []string{"John", "Jane"})
}

func CreateUser(c echo.Context) error {
	user := new(struct {
		Name string `json:"name"`
	})
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	return c.JSON(http.StatusCreated, user)
}
