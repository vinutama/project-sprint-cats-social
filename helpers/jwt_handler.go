package helpers

import (
	cfg "cats-social/config"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func GetTokenHandler() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(cfg.EnvConfigs.JwtSecret)},
		ContextKey: JwtContextKey,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return fiber.NewError(fiber.StatusForbidden, err.Error())
		},
	})
}

func CheckTokenHeader(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	} else {
		return ctx.Next()
	}
}
