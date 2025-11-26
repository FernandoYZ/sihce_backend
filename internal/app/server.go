package app

import (
	"github.com/gofiber/fiber/v2"
	"backend/internal/config"
	"backend/internal/config/database"
)

type App struct {
	Fiber  	*fiber.App
	Config 	*config.Config
	Db		*database.GestorDB
}

func New(cfg *config.Config, db *database.GestorDB) *App {
	app := fiber.New(fiber.Config{
		AppName:      "Backend API v1.0",
		ErrorHandler: ErroresGlobales,
	})

	ConfigurarMiddlewares(app, cfg)
	ConfigurarRutas(app, db)

	return &App{
		Fiber:  app,
		Config: cfg,
		Db:     db,
	}
}