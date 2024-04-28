package auth_service

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

type AuthService interface {
	GenerateToken(ctx context.Context, userId string) (string, error)
	GetValidUser(ctx *fiber.Ctx) (string, error)
}
