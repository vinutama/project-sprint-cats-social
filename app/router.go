package app

import (
	"cats-social/controller"
	"cats-social/helpers"
	cat_repository "cats-social/repository/cat"
	match_repository "cats-social/repository/match"
	user_repository "cats-social/repository/user"
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

	userRepository := user_repository.NewUserRepository(dbPool)
	userService := user_service.NewUserService(userRepository, validator)
	userController := controller.NewUserController(userService)

	catRepository := cat_repository.NewCatRepository(dbPool)
	catService := cat_service.NewCatService(catRepository, validator)
	catController := controller.NewCatController(catService)

	matchRepository := match_repository.NewMatchRepository(dbPool)
	matchService := match_service.NewMatchService(matchRepository, validator)
	matchController := controller.NewMatchController(matchService)

	// Users API
	userApi := app.Group("/v1/user")
	userApi.Post("/register", userController.Register)
	userApi.Post("/login", userController.Login)

	// JWT middleware
	// app.Use(helpers.CheckTokenHeader)
	app.Use(helpers.GetTokenHandler())

	// All request below here shoud use Bearer <token>

	// Cats API
	catApi := app.Group("/v1/cat")
	catApi.Post("/", catController.Create)
	catApi.Put("/:id", catController.EditCat)
	catApi.Get("/", catController.Search)
	catApi.Delete("/:id", catController.Delete)

	// Match API
	matchApi := catApi.Group("/match")
	matchApi.Post("/", matchController.Create)
	matchApi.Get("/", matchController.Get)
	matchApi.Post("/approve", matchController.Approve)
	matchApi.Post("/reject", matchController.Reject)
	matchApi.Delete("/:id", matchController.Delete)
}
