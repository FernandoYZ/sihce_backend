package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"backend/internal/config"
	"backend/internal/config/database"
	"backend/internal/modules/triaje"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Cargar configuración general
	cfg, err := config.Cargar("config.yml")
	if err != nil {
		log.Fatalf("Error al cargar configuración: %v", err)
	}

	// Inicializar gestor de base de datos
	gestor, err := database.NuevoGestor("internal/config/database/config.yml")
	if err != nil {
		log.Fatalf("Error al inicializar gestor de BD: %v", err)
	}
	defer gestor.Cerrar()

	// Verificar conexión inicial
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := gestor.VerificarSalud(ctx); err != nil {
		log.Printf("Advertencia BD: %v", err)
	}

	// Crear aplicación Fiber
	app := fiber.New(fiber.Config{
		AppName:      "Backend API v1.0",
		ErrorHandler: manejadorErrores,
	})

	// Middlewares
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} ${method} ${path} ${latency}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: strings.Join(cfg.App.CorsOrigins, ","),
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Rutas
	app.Get("/health", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
		defer cancel()

		if err := gestor.VerificarSalud(ctx); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "unhealthy",
				"error":  err.Error(),
			})
		}
		return c.JSON(fiber.Map{"status": "healthy"})
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"mensaje": "API funcionando correctamente",
			"version": "1.0",
		})
	})

	// Registrar módulos
	api := app.Group("/api")
	registrarModulos(api, gestor)

	// Graceful shutdown
	canal := make(chan os.Signal, 1)
	signal.Notify(canal, os.Interrupt, syscall.SIGTERM)

	go func() {
		puerto := fmt.Sprintf(":%d", cfg.App.Port)
		log.Printf("Servidor iniciado en puerto %d", cfg.App.Port)
		if err := app.Listen(puerto); err != nil {
			log.Fatalf("Error al iniciar servidor: %v", err)
		}
	}()

	<-canal
	log.Println("Cerrando servidor...")
	app.ShutdownWithTimeout(10 * time.Second)
}

func manejadorErrores(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error":   true,
		"message": err.Error(),
	})
}

func registrarModulos(api fiber.Router, gestor *database.GestorDB) {
	// Cada módulo se inicializa de forma independiente
	// Si un módulo falla, no afecta a los demás

	// Módulo Triaje
	moduloTriaje := triaje.NuevoModulo(gestor)
	moduloTriaje.RegistrarRutas(api)

	// Aquí se registrarán más módulos...
}
