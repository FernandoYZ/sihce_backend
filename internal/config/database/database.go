package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/microsoft/go-mssqldb"
	"gopkg.in/yaml.v3"
)

// GestorDB gestiona los pools de conexión a las bases de datos
type GestorDB struct {
	principal     *sql.DB
	secundaria    *sql.DB
	configuracion *Configuracion
	mu            sync.RWMutex
}

var (
	instancia *GestorDB
	once      sync.Once
	mu        sync.Mutex
)

// NuevoGestor crea y retorna la instancia singleton del gestor de bases de datos
func NuevoGestor(rutaConfig string) (*GestorDB, error) {
	var err error

	once.Do(func() {
		var config *Configuracion
		config, err = cargarConfiguracion(rutaConfig)
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

// cargarConfiguracion carga la configuración desde el archivo YAML
func cargarConfiguracion(rutaConfig string) (*Configuracion, error) {
	rutaConfig = os.ExpandEnv(rutaConfig)

	archivo, err := os.ReadFile(rutaConfig)
	if err != nil {
		return nil, fmt.Errorf("error al leer archivo de configuración: %w", err)
	}

	contenidoExpandido := os.ExpandEnv(string(archivo))

	var config Configuracion
	if err := yaml.Unmarshal([]byte(contenidoExpandido), &config); err != nil {
		return nil, fmt.Errorf("error al parsear YAML: %w", err)
	}

	return &config, nil
}

// construirCadenaConexion construye la cadena de conexión para SQL Server
func construirCadenaConexion(cfg ConfiguracionDB) string {
	encrypt := "disable"
	if cfg.Encrypt {
		encrypt = "true"
	}

	trustCert := "false"
	if cfg.TrustServerCert {
		trustCert = "true"
	}

	return fmt.Sprintf(
		"server=%s;port=%d;database=%s;user id=%s;password=%s;encrypt=%s;TrustServerCertificate=%s;connection timeout=%d",
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.User,
		cfg.Password,
		encrypt,
		trustCert,
		cfg.Pool.ConnectionTimeoutMs/1000,
	)
}

// inicializarPool inicializa un pool de conexiones
func inicializarPool(cfg ConfiguracionDB, nombre string) (*sql.DB, error) {
	cadenaConexion := construirCadenaConexion(cfg)

	db, err := sql.Open("sqlserver", cadenaConexion)
	if err != nil {
		return nil, fmt.Errorf("error al abrir conexión %s: %w", nombre, err)
	}

	// Configurar el pool de conexiones
	db.SetMaxOpenConns(cfg.Pool.Max)
	db.SetMaxIdleConns(cfg.Pool.Min)
	db.SetConnMaxIdleTime(time.Duration(cfg.Pool.IdleTimeoutMs) * time.Millisecond)
	db.SetConnMaxLifetime(30 * time.Minute)

	// Verificar la conexión con contexto y timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Pool.ConnectionTimeoutMs)*time.Millisecond)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("error al hacer ping a la base de datos %s: %w", nombre, err)
	}

	log.Printf("[Database] Pool de conexión %s inicializado correctamente (DB: %s)", nombre, cfg.Name)
	return db, nil
}

// ObtenerPrincipal obtiene el pool de conexión principal (lazy initialization)
func (g *GestorDB) ObtenerPrincipal() (*sql.DB, error) {
	if g.principal != nil {
		return g.principal, nil
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	// Double-check locking pattern
	if g.principal != nil {
		return g.principal, nil
	}

	db, err := inicializarPool(g.configuracion.Database.Principal, "principal")
	if err != nil {
		return nil, err
	}

	g.principal = db
	return g.principal, nil
}

// ObtenerSecundaria obtiene el pool de conexión secundaria (lazy initialization)
func (g *GestorDB) ObtenerSecundaria() (*sql.DB, error) {
	if g.secundaria != nil {
		return g.secundaria, nil
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	// Double-check locking pattern
	if g.secundaria != nil {
		return g.secundaria, nil
	}

	db, err := inicializarPool(g.configuracion.Database.Secundaria, "secundaria")
	if err != nil {
		return nil, err
	}

	g.secundaria = db
	return g.secundaria, nil
}

// VerificarSalud verifica el estado de las conexiones activas
func (g *GestorDB) VerificarSalud(ctx context.Context) error {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if g.principal != nil {
		if err := g.principal.PingContext(ctx); err != nil {
			return fmt.Errorf("error en health check de base de datos principal: %w", err)
		}
	}

	if g.secundaria != nil {
		if err := g.secundaria.PingContext(ctx); err != nil {
			return fmt.Errorf("error en health check de base de datos secundaria: %w", err)
		}
	}

	return nil
}

// Cerrar cierra todas las conexiones de base de datos
func (g *GestorDB) Cerrar() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	var errores []error

	if g.principal != nil {
		if err := g.principal.Close(); err != nil {
			errores = append(errores, fmt.Errorf("error al cerrar pool principal: %w", err))
		} else {
			log.Println("[Database] Pool de conexión principal cerrado correctamente")
		}
		g.principal = nil
	}

	if g.secundaria != nil {
		if err := g.secundaria.Close(); err != nil {
			errores = append(errores, fmt.Errorf("error al cerrar pool secundaria: %w", err))
		} else {
			log.Println("[Database] Pool de conexión secundaria cerrado correctamente")
		}
		g.secundaria = nil
	}

	if len(errores) > 0 {
		return fmt.Errorf("errores al cerrar conexiones: %v", errores)
	}

	return nil
}
