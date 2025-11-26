package database

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/microsoft/go-mssqldb"
)

type GestorDB struct {
	principal     *sql.DB
	secundaria    *sql.DB
	configuracion *Configuracion
	mu            sync.RWMutex
}

var (
	instancia *GestorDB
	once      sync.Once
)

// NuevoGestor crea y retorna la instancia singleton del gestor de bases de datos
func NuevoGestor(rutaConfig string) (*GestorDB, error) {
	var err error

	once.Do(func() {
		var config *Configuracion
		config, err = CargarConfiguracion(rutaConfig)
		if err != nil {
			err = fmt.Errorf("error al cargar configuración: %w", err)
			return
		}

		instancia = &GestorDB{
			configuracion: config,
		}

		log.Println("[Database] Gestor de base de datos inicializado")
	})

	if err != nil {
		return nil, err
	}

	return instancia, nil
}

// ObtenerGestor devuelve la instancia singleton del gestor
// Retorna error si no ha sido inicializado con NuevoGestor
func ObtenerGestor() (*GestorDB, error) {
	if instancia == nil {
		return nil, fmt.Errorf("gestor no inicializado, llamar NuevoGestor primero")
	}
	return instancia, nil
}
