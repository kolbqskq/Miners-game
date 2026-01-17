package main

import (
	"encoding/gob"
	"miners_game/config"
	"miners_game/internal/auth"
	"miners_game/internal/game"
	"miners_game/internal/game/loop"
	"miners_game/internal/game/sessions"
	"miners_game/internal/pages"
	"miners_game/internal/user"
	"miners_game/pkg/database"
	"miners_game/pkg/logger"
	"miners_game/pkg/middleware"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/postgres/v3"
	"github.com/gookit/validate/locales/ruru"
	"github.com/rs/zerolog"
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
	loopService := loop.NewService()
	sessionService := sessions.NewService(time.Minute*5)
	gameService := game.NewService(game.ServiceDeps{
		Repo:     gameRepository,
		Loop:     loopService,
		Sessions: sessionService,
	})
	authService := auth.NewService(auth.ServiceDeps{
		UserRepository: userRepository,
		GmailConfig:    gmailConfig,
	})

	//Handlers:
	pages.NewHandler(pages.HandlerDeps{
		Router: app,
		Logger: customLogger,
		Store: store,
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

	App(loopService, gameService, customLogger)

	if err := app.Listen(":3000"); err != nil {
		customLogger.Fatal().Err(err).Msg("не удалось запустить HTTP сервер")
	}
}

func App(loopService *loop.Service, gameService *game.Service, logger *zerolog.Logger) {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			now := time.Now().Unix()
			loopService.Tick(now)
		}
	}()

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			gameService.HandleExpiredSessions()
		}
	}()

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			gameService.SaveAll()
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(
		sigCh,
		os.Interrupt,
		syscall.SIGTERM,
	)
	go func() {
		<-sigCh
		gameService.SaveAll()
		logger.Info().Msg("games saves complete")
		os.Exit(0)
	}()

}
