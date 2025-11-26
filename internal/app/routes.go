package app

import (
	"backend/internal/config/database"
	"backend/internal/modules/triaje"

	"github.com/gofiber/fiber/v2"
)

func ConfigurarRutas(router *fiber.App, db *database.GestorDB) {
	api := router.Group("/api")
	api.Get("/", VerificarApi(db))

	// Registro de módulos de la API
	triaje.NuevoModulo(db).RegistrarRutas(api)
}
