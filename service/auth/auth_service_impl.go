package auth_service

import (
	cfg "cats-social/config"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var expDuration = time.Now().Add(time.Minute * 30).Unix()

type AuthServiceImpl struct {
}

func NewAuthService() AuthService {
	return &AuthServiceImpl{}
}

func (service *AuthServiceImpl) GenerateToken(ctx context.Context, userId string) (string, error) {
	jwtconf := jwt.MapClaims{
		"user_id": userId,
		"exp":     expDuration,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtconf)
	signToken, err := token.SignedString([]byte(cfg.EnvConfigs.JwtSecret))
	if err != nil {
		return "", err
	}

	return signToken, nil
}

func (service *AuthServiceImpl) GetValidUser(ctx *fiber.Ctx) (string, error) {
	userInfo := ctx.Locals("userInfo").(*jwt.Token)
	// convert userInfo claims to jwt mapclaims
	jwtconf := userInfo.Claims.((jwt.MapClaims))
	// convert user_id to string
	userId := jwtconf["user_id"].(string)

	return userId, nil
}
