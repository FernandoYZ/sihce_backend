package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/microsoft/go-mssqldb"
)

// Variables globales para los pools de conexión
var (
	PoolPrincipal *sql.DB
	PoolExterno   *sql.DB
	mutexPrincipal sync.Mutex
	mutexExterno   sync.Mutex
)

// ConfigDB estructura para la configuración de la base de datos
type ConfigDB struct {
	Usuario              string
	Password             string
	Servidor             string
	BaseDatos            string
	Puerto               string
	Encrypt              string
	TrustServerCert      string
}

// obtenerConfigPrincipal obtiene la configuración para la base de datos principal
func obtenerConfigPrincipal() ConfigDB {
	return ConfigDB{
		Usuario:         os.Getenv("DB_USER"),
		Password:        os.Getenv("DB_PASSWORD"),
		Servidor:        os.Getenv("DB_SERVER"),
		BaseDatos:       os.Getenv("DB_DATABASE_PRINCIPAL"),
		Puerto:          os.Getenv("DB_PORT"),
		Encrypt:         os.Getenv("DB_ENCRYPT"),
		TrustServerCert: os.Getenv("DB_TRUST_SERVER_CERTIFICATE"),
	}
}

// obtenerConfigExterna obtiene la configuración para la base de datos externa
func obtenerConfigExterna() ConfigDB {
	return ConfigDB{
		Usuario:         os.Getenv("DB_USER"),
		Password:        os.Getenv("DB_PASSWORD"),
		Servidor:        os.Getenv("DB_SERVER"),
		BaseDatos:       os.Getenv("DB_DATABASE_SECUNDARIA"),
		Puerto:          os.Getenv("DB_PORT"),
		Encrypt:         os.Getenv("DB_ENCRYPT"),
		TrustServerCert: os.Getenv("DB_TRUST_SERVER_CERTIFICATE"),
	}
}

// construirConnectionString construye la cadena de conexión para SQL Server
func construirConnectionString(config ConfigDB) string {
	return fmt.Sprintf(
		"server=%s;user id=%s;password=%s;port=%s;database=%s;encrypt=disable;TrustServerCertificate=true",
		config.Servidor,
		config.Usuario,
		config.Password,
		config.Puerto,
		config.BaseDatos,
	)
}

// InicializarPoolPrincipal inicializa el pool de conexión principal
func InicializarPoolPrincipal() error {
	mutexPrincipal.Lock()
	defer mutexPrincipal.Unlock()

	if PoolPrincipal != nil {
		return nil
	}

	config := obtenerConfigPrincipal()
	connString := construirConnectionString(config)

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return fmt.Errorf("error al abrir conexión principal: %w", err)
	}

	// Configurar el pool de conexiones
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	// Verificar la conexión
	if err := db.Ping(); err != nil {
		db.Close()
		return fmt.Errorf("error al hacer ping a la base de datos principal: %w", err)
	}

	PoolPrincipal = db
	log.Printf("[Config] Pool de conexión principal inicializado (%s)", config.BaseDatos)
	return nil
}

// InicializarPoolExterno inicializa el pool de conexión externo
func InicializarPoolExterno() error {
	mutexExterno.Lock()
	defer mutexExterno.Unlock()

	if PoolExterno != nil {
		return nil
	}

	config := obtenerConfigExterna()
	connString := construirConnectionString(config)

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return fmt.Errorf("error al abrir conexión externa: %w", err)
	}

	// Configurar el pool de conexiones
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	// Verificar la conexión
	if err := db.Ping(); err != nil {
		db.Close()
		return fmt.Errorf("error al hacer ping a la base de datos externa: %w", err)
	}

	PoolExterno = db
	log.Printf("[Config] Pool de conexión externa inicializado (%s)", config.BaseDatos)
	return nil
}

// ObtenerConexionPrincipal obtiene el pool de conexión principal
func ObtenerConexionPrincipal() (*sql.DB, error) {
	if PoolPrincipal == nil {
		if err := InicializarPoolPrincipal(); err != nil {
			return nil, err
		}
	}
	return PoolPrincipal, nil
}

// ObtenerConexionExterna obtiene el pool de conexión externa
func ObtenerConexionExterna() (*sql.DB, error) {
	if PoolExterno == nil {
		if err := InicializarPoolExterno(); err != nil {
			return nil, err
		}
	}
	return PoolExterno, nil
}

// CerrarConexiones cierra todos los pools de conexión
func CerrarConexiones() {
	if PoolPrincipal != nil {
		if err := PoolPrincipal.Close(); err != nil {
			log.Printf("[Config] Error al cerrar pool principal: %v", err)
		} else {
			log.Println("[Config] Pool de conexión principal cerrado")
		}
		PoolPrincipal = nil
	}

	if PoolExterno != nil {
		if err := PoolExterno.Close(); err != nil {
			log.Printf("[Config] Error al cerrar pool externo: %v", err)
		} else {
			log.Println("[Config] Pool de conexión externa cerrado")
		}
		PoolExterno = nil
	}
}
