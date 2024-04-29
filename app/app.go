package app

import (
	cfg "cats-social/config"
	"cats-social/database"
	"cats-social/script"
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
	// Temporary helper to initiate tables
	if err := script.InitiateTables(dbPool); err != nil {
		log.Fatal("Error when initializing tables:", err)
	}
	defer dbPool.Close()

	app.Use(logger.New())
	// Register BP
	RegisterBluePrint(app, dbPool)

	err := app.Listen(":8080")
	log.Fatal(err)
}
