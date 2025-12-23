package main

import (
	"miners_game/config"
	"miners_game/internal/home"
	"miners_game/pkg/database"
	"miners_game/pkg/logger"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	//config
	config.Init()
	loggerConfig := config.NewLogConfig()
	dbConfig := config.NewDatabaseConfig()

	customLogger := logger.NewLogger(loggerConfig)

	app := fiber.New()

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: customLogger,
	}))
	app.Use(recover.New())
	app.Static("/public", "./public")
	dbPool := database.CreateDbPool(dbConfig, customLogger)
	defer dbPool.Close()

	//Handlers:
	home.NewHandler(home.HandlerDeps{
		Router: app,
		Logger: customLogger,
	})

	if err := app.Listen(":3000"); err != nil {
		customLogger.Fatal().Err(err).Msg("не удалось запустить HTTP сервер")
	}
}
