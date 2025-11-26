package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

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

// InicializarPrincipal fuerza la conexión a la BD principal para verificar credenciales al arranque
func (g *GestorDB) InicializarPrincipal() error {
	_, err := g.ObtenerPrincipal()
	return err
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

// construirCadenaConexion construye la cadena de conexión para SQL Server
func construirCadenaConexion(cfg ConfiguracionDB) string {
	return fmt.Sprintf(
		"server=%s;port=%d;database=%s;user id=%s;password=%s;encrypt=%t;TrustServerCertificate=%t;connection timeout=%d",
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.User,
		cfg.Password,
		cfg.Encrypt,
		cfg.TrustServerCert,
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

	db.SetMaxOpenConns(cfg.Pool.Max)
	db.SetMaxIdleConns(cfg.Pool.Min)
	db.SetConnMaxIdleTime(time.Duration(cfg.Pool.IdleTimeoutMs) * time.Millisecond)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Pool.ConnectionTimeoutMs)*time.Millisecond)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("error al hacer ping a la base de datos %s: %w", nombre, err)
	}

	log.Printf("[Database] Pool de conexión %s inicializado correctamente (DB: %s)", nombre, cfg.Name)
	return db, nil
}
