package utils

import (
	"backend_tech_movement_hex/internal/core/domain"

	"github.com/go-playground/validator"
)

var validate = validator.New()

func ValidateNewsInput(input *domain.NewsRequest) error {
	if err := validate.Struct(input); err != nil {
		return err
	}
	return nil
}
