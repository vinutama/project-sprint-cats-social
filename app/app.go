package app

import (
	cfg "cats-social/config"
	"log"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func StartApp() {
	app := fiber.New(fiber.Config{
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		IdleTimeout:  cfg.IdleTimeout,
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		Prefork:      cfg.Prefork,
	})

	app.Use(logger.New())

	// Register BP
	RegisterBluePrint(app)

	err := app.Listen("localhost:8000")
	log.Fatal(err)
}
