package app

import (
	"fmt"
	"log"

	"backend/internal/config"
	"backend/internal/config/database"

	"github.com/gofiber/fiber/v2"
)

type App struct {
	Fiber  *fiber.App
	Config *config.Config
	Db     *database.GestorDB
}

func New(cfg *config.Config, db *database.GestorDB) *App {
	app := fiber.New(fiber.Config{
		AppName:               "API Consulta Externa v3.0 - Backend",
		ErrorHandler:          ErroresGlobales,
		DisableStartupMessage: false,
	})

	ConfigurarMiddlewares(app, cfg)
	ConfigurarRutas(app, db)

	return &App{
		Fiber:  app,
		Config: cfg,
		Db:     db,
	}
}

// Run inicia el servidor escuchando en el puerto configurado
func (a *App) Run() error {
	puerto := fmt.Sprintf(":%d", a.Config.App.Port)
	return a.Fiber.Listen(puerto)
}


func (a *App) Shutdown() error {
	log.Println("🛑 Apagando servidor HTTP...")
	return a.Fiber.Shutdown()
}