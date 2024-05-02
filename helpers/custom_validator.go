package helpers

import (
	"regexp"
	"strconv"

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

func validateBoolean(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	_, err := strconv.ParseBool(value)

	return err == nil
}

func validateRegex(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	regexTag := fl.Param()
	if regexTag == "" {
		return false
	}

	matched, _ := regexp.MatchString(regexTag, value)
	return matched
}

func RegisterCustomValidator(validator *validator.Validate) {
	// validator.RegisterValidation() -> if you want to create new tags rule to be used on struct entity
	// validator.RegisterStructValidation() -> if you want to create validator then access all fields to the struct entity

	validator.RegisterValidation("sex", validateGender)
	validator.RegisterValidation("catRace", validateCatRace)
	validator.RegisterValidation("boolean", validateBoolean)
	validator.RegisterValidation("regex", validateRegex)
}
