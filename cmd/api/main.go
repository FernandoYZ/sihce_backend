package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"backend/internal/app"
	"backend/internal/config"
	"backend/internal/config/database"
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

	// if err := gestor.InicializarPrincipal(); err != nil {
	// 	log.Fatalf("Error al inicializar la conexión con la base de datos principal: %v", err)
	// }

	// Crear una nueva instancia de la aplicación, delegando toda la configuración.
	servidor := app.New(cfg, gestor)

	// Cerrar el servidor al recibir una señal de terminación.
	detener := make(chan os.Signal, 1)
	signal.Notify(detener, os.Interrupt, syscall.SIGTERM)

	// Iniciar el servidor en una goroutine para no bloquear el canal de apagado.
	go func() {
		if err := servidor.Run(); err != nil {
			log.Fatalf("Error al iniciar servidor: %v", err)
		}
	}()

	// Bloquear hasta que se reciba una señal de apagado.
	<-detener
	if err := servidor.Shutdown(); err != nil {
		log.Printf("Error durante el apagado del servidor: %v", err)
	}
	log.Println("Servidor cerrado correctamente.")
}
