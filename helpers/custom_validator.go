package helpers

import (
	"github.com/go-playground/validator"
)

func validateGender(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return value == "male" || value == "female"
}

func validateCatRace(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	for _, race := range CatRaces {
		if race == value {
			return true
		}
	}
	return false
}

func RegisterCustomValidator(validator *validator.Validate) {
	// validator.RegisterValidation() -> if you want to create new tags rule to be used on struct entity
	// validator.RegisterStructValidation() -> if you want to create validator then access all fields to the struct entity

	validator.RegisterValidation("sex", validateGender)
	validator.RegisterValidation("catRace", validateCatRace)
}
