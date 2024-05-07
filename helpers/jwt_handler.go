package helpers

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

func GetTokenHandler() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(viper.GetString("JWT_SECRET"))},
		ContextKey: JwtContextKey,
		SuccessHandler: func(c *fiber.Ctx) error {
			auth := c.Locals(JwtContextKey).(*jwt.Token)
			claims := auth.Claims.(jwt.MapClaims)
			c.Locals("userId", claims["user_id"].(string))
			return c.Next()
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}})
}

func CheckTokenHeader(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	} else {
		return ctx.Next()
	}
}
