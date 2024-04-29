package helpers

import "github.com/gofiber/fiber/v2"

var ErrorBadRequest = fiber.NewError(fiber.StatusBadRequest, "Bad request")
