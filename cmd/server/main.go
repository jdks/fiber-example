package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/rs/zerolog"

	"github.com/jdks/fiber-example/internal/config"
	"github.com/jdks/fiber-example/internal/server"
)

func setupRoutes(app *fiber.App, s server.Server) {
	app.Get("/users", s.GetAllUsers)
	app.Get("/users/:user_id", s.GetUser)
	app.Get("/events", s.QueryEvents)
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log := zerolog.New(os.Stdout).With().Logger()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load configs")
	}

	app := fiber.New(cfg.Fiber)
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{Output: log}))

	s, err := server.New(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create server")
	}
	setupRoutes(app, s)

	app.Listen(fmt.Sprintf("%s:%s", cfg.App.Host, cfg.App.Port))
}
