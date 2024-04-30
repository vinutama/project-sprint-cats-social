package exceptions

import (
	"github.com/gofiber/fiber/v2"
)

// http status code https://pkg.go.dev/github.com/gofiber/fiber/v2@v2.52.4#StatusConflict

func Exception(ctx *fiber.Ctx, err error) error {
	if fiberErr, ok := err.(*fiber.Error); ok {
		// If it's a Fiber error, return it as JSON
		return ctx.Status(fiberErr.Code).JSON(fiber.Map{"message": fiberErr.Message})
	}
	// Otherwise, return the error as is
	return err
}

func ConflictException(message string) error {
	return fiber.NewError(fiber.StatusConflict, message)
}

func NotFoundException(message string) error {
	return fiber.NewError(fiber.StatusNotFound, message)
}

func BadRequestException(message string) error {
	return fiber.NewError(fiber.StatusBadRequest, message)
}

func InternalServerException(message string) error {
	return fiber.NewError(fiber.StatusInternalServerError, message)
}

func UnauthorizedException(message string) error {
	return fiber.NewError(fiber.StatusUnauthorized, message)
}
