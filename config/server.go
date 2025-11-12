package config

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

// NuevaApp crea una nueva instancia de Fiber
func NuevaApp() *fiber.App {
	return fiber.New()
}

// ObtenerPuerto obtiene el puerto desde las variables de entorno
func ObtenerPuerto() string {
	puerto := os.Getenv("PORT")
	if puerto == "" {
		puerto = "3054"
	}
	return ":" + puerto
}
