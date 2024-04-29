package app

import (
	"cats-social/controller"
	"cats-social/helpers"
	user_repository "cats-social/repository/user"
	auth_service "cats-social/service/auth"
	user_service "cats-social/service/user"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterBluePrint(app *fiber.App, dbPool *pgxpool.Pool) {
	//TODO:  add validator Here

	authService := auth_service.NewAuthService()

	userRepository := user_repository.NewUserRepository()
	userService := user_service.NewUserService(userRepository, dbPool, authService)
	userController := controller.NewUserController(userService, authService)

	app.Post("/v1/user/register", userController.Register)
	app.Post("/v1/user/login", userController.Login)

	app.Get("/v1/user/hehe", func(c *fiber.Ctx) error {
		return c.SendString("APANIH HEHE")
	})

	//middleware JWT down after this code
	app.Use(helpers.CheckTokenHeader)
	app.Use(helpers.GetTokenHandler())
}
