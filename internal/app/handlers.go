package app

import (
	"backend/internal/config"
	"backend/internal/config/database"
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Función para manejo de errores globales
func ErroresGlobales(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	mensaje := "Ha ocurrido un error interno en el servidor."
	tipo := "INTERNAL_SERVER_ERROR"

	// Si el error es de Fiber, usar su código
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		if code == fiber.StatusNotFound {
			mensaje = "Recurso no encontrado."
			tipo = "NOT_FOUND"
		} else if e.Message != "" {
			mensaje = e.Message
		}
	}

	// Log interno (solo para servidor)
	if code == fiber.StatusInternalServerError {
		log.Printf(
			"[ERROR] %s %s -> %v",
			c.Method(),
			c.OriginalURL(),
			err,
		)
	}

	// Crear el mapa de respuesta base
	respuesta := fiber.Map{
		"error":     true,
		"tipo":      tipo,
		"mensaje":   mensaje,
		"path":      c.OriginalURL(),
		"metodo":    c.Method(),
		"status":    code,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	// Añadir detalles del error solo en entorno de desarrollo para errores 500
	cfg := config.Obtener()
	if cfg != nil && cfg.App.AppEnv == "dev" && code == fiber.StatusInternalServerError {
		respuesta["detalles"] = err.Error()
	}

	// Respuesta estandarizada
	return c.Status(code).JSON(respuesta)
}

// VerificarApiHandler crea un handler para verificar el estado de la API y sus dependencias.
func VerificarApi(db *database.GestorDB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Verificar salud de la BD principal
		ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
		defer cancel()

		ESTADO_DB_PRINCIPAL := true
		if err := db.VerificarSalud(ctx); err != nil {
			// Aquí podrías loggear el error `err` para depuración
			ESTADO_DB_PRINCIPAL = false
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": true,
			"data": fiber.Map{
				"mensaje":   "API funcionando correctamente",
				"version":   "v3.0",
				"timestamp": time.Now().UTC().Format(time.RFC3339),
				"servicios": fiber.Map{
					"database": fiber.Map{
						"DB_SIGH":         ESTADO_DB_PRINCIPAL,
						"DB_SIGH_EXTERNA": "lazy",
					},
				},
			},
		})
	}
}
