package app

import (
	cfg "cats-social/config"
	"cats-social/database"
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

	dbPool := database.GetConnPool()
	defer dbPool.Close()

	app.Use(logger.New())
	// Register BP
	RegisterBluePrint(app, dbPool)

	err := app.Listen("localhost:8080")
	log.Fatal(err)
}
