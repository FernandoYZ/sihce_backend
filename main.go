package main

import (
	"backend/config"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No se pudo cargar .env:", err)
	}

	// Conectar a las bases de datos
	if err := config.InicializarPoolPrincipal(); err != nil {
		log.Printf("ERROR DB Principal: %v", err)
	}

	if err := config.InicializarPoolExterno(); err != nil {
		log.Printf("ERROR DB Externa: %v", err)
	}

	// Crear app Fiber
	app := config.NuevaApp()

	// Ruta de prueba
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"mensaje": "Servidor funcionando",
		})
	})

	// Iniciar servidor
	puerto := config.ObtenerPuerto()
	fmt.Printf("Servidor corriendo en http://localhost%s\n", puerto)
	log.Fatal(app.Listen(puerto))
}
