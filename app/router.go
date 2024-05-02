package app

import (
	"cats-social/controller"
	"cats-social/helpers"
	cat_repository "cats-social/repository/cat"
	match_repository "cats-social/repository/match"
	user_repository "cats-social/repository/user"
	auth_service "cats-social/service/auth"
	cat_service "cats-social/service/cat"
	match_service "cats-social/service/match"
	user_service "cats-social/service/user"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterBluePrint(app *fiber.App, dbPool *pgxpool.Pool) {
	validator := validator.New()
	// register custom validator
	helpers.RegisterCustomValidator(validator)

	authService := auth_service.NewAuthService()

	userRepository := user_repository.NewUserRepository()
	userService := user_service.NewUserService(userRepository, dbPool, authService, validator)
	userController := controller.NewUserController(userService, authService)

	catRepository := cat_repository.NewCatRepository()
	catService := cat_service.NewCatService(catRepository, dbPool, authService, validator)
	catController := controller.NewCatController(catService, authService)

	matchRepository := match_repository.NewMatchRepository()
	matchService := match_service.NewMatchService(matchRepository, dbPool, authService, validator)
	matchController := controller.NewMatchController(matchService, authService)

	// Users API
	userApi := app.Group("/v1/user")
	userApi.Post("/register", userController.Register)
	userApi.Post("/login", userController.Login)

	// JWT middleware
	app.Use(helpers.CheckTokenHeader)
	app.Use(helpers.GetTokenHandler())

	// All request below here shoud use Bearer <token>

	// Cats API
	catApi := app.Group("/v1/cat")
	catApi.Post("/", catController.Create)
	catApi.Get("/", catController.Search)

	// Match API
	matchApi := catApi.Group("/match")
	matchApi.Post("/", matchController.Create)
}
