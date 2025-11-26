package app

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// Función para manejo de errores globales
func ErroresGlobales(h *fiber.Ctx, err error) error {
	return h.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"status":  false,
		"error":   "INTERNAL_SERVER_ERROR",
		"mensaje": err.Error(),
	})
}

// Función para verificar el estado de la API
func VerificarApi(h *fiber.Ctx) error {
	ESTADO_SIGH := true
	ESTADO_SIGH_EXTERNA := true
	ESTADO_SOCKET := false
	ESTADO_TEMPLATES := false

	return h.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"data": fiber.Map{
			"mensaje":  "API funcionando correctamente",
			"version":  "v3.0",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"servicios": fiber.Map{
				"base_datos": fiber.Map{
					"SIGH":        ESTADO_SIGH,
					"SIGH_EXTERNA": ESTADO_SIGH_EXTERNA,
				},
				"WebSocket": ESTADO_SOCKET,
				"Templates": ESTADO_TEMPLATES,
			},
			"modulos": fiber.Map{
				"triaje":    ESTADO_TEMPLATES,
				"atenciones": ESTADO_TEMPLATES,
			},
		},
	})
}
