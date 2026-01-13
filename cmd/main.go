package main

import (
	"encoding/gob"
	"miners_game/config"
	"miners_game/internal/auth"
	"miners_game/internal/game"
	"miners_game/internal/pages"
	"miners_game/internal/user"
	"miners_game/pkg/database"
	"miners_game/pkg/logger"
	"miners_game/pkg/middleware"
	"time"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/postgres/v3"
	"github.com/gookit/validate/locales/ruru"
)

func main() {
	//configs
	config.Init()
	loggerConfig := config.NewLogConfig()
	dbConfig := config.NewDatabaseConfig()
	gmailConfig := config.NewGmailConfig()
	ruru.RegisterGlobal()

	customLogger := logger.NewLogger(loggerConfig)

	app := fiber.New()

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: customLogger,
	}))
	app.Use(recover.New())
	app.Static("/public", "./public")
	dbPool := database.CreateDbPool(dbConfig, customLogger)
	defer dbPool.Close()
	storage := postgres.New(postgres.Config{
		DB:         dbPool,
		Table:      "session",
		Reset:      false,
		GCInterval: 10 * time.Second,
	})
	store := session.New(session.Config{
		Storage: storage,
	})
	gob.Register(auth.RegisterSession{})
	app.Use(middleware.AuthMiddleware(store))

	//Repositories:
	gameRepository := game.NewRepository(game.RepositoryDeps{
		DbPool: dbPool,
		Logger: customLogger,
	})
	userRepository := user.NewRepository(user.RepositoryDeps{
		DbPool: dbPool,
		Logger: customLogger,
	})
	//Services:
	gameService := game.NewService(game.ServiceDeps{
		GameRepository: gameRepository,
	})
	authService := auth.NewService(auth.ServiceDeps{
		UserRepository: userRepository,
		GmailConfig:    gmailConfig,
	})

	//Handlers:
	pages.NewHandler(pages.HandlerDeps{
		Router: app,
		Logger: customLogger,
	})
	game.NewHandler(game.HandlerDeps{
		Router:      app,
		Logger:      customLogger,
		GameService: gameService,
		Store:       store,
	})
	auth.NewHandler(auth.HandlerDeps{
		Router:      app,
		Logger:      customLogger,
		AuthService: authService,
		Store:       store,
	})

	if err := app.Listen(":3000"); err != nil {
		customLogger.Fatal().Err(err).Msg("не удалось запустить HTTP сервер")
	}
}
