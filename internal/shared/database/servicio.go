package database

import (
	"context"
	"database/sql"
	"fmt"

	"backend/internal/config/database"
)

type ServicioDB struct {
	gestor *database.GestorDB
}

func NuevoServicio(gestor *database.GestorDB) *ServicioDB {
	return &ServicioDB{gestor: gestor}
}

// ObtenerConexion obtiene la conexión según el parámetro usarSecundaria
func (s *ServicioDB) ObtenerConexion(usarSecundaria bool) (*sql.DB, error) {
	if usarSecundaria {
		return s.gestor.ObtenerSecundaria()
	}
	return s.gestor.ObtenerPrincipal()
}

// EjecutarQuery ejecuta un query SQL
// Por defecto usa la BD principal, si usarSecundaria=true usa la secundaria
func (s *ServicioDB) EjecutarQuery(ctx context.Context, query string, usarSecundaria bool, args ...interface{}) (*sql.Rows, error) {
	db, err := s.ObtenerConexion(usarSecundaria)
	if err != nil {
		return nil, fmt.Errorf("error al obtener conexión: %w", err)
	}

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar query: %w", err)
	}

	return rows, nil
}

// EjecutarQueryRow ejecuta un query que retorna una sola fila
func (s *ServicioDB) EjecutarQueryRow(ctx context.Context, query string, usarSecundaria bool, args ...interface{}) *sql.Row {
	db, err := s.ObtenerConexion(usarSecundaria)
	if err != nil {
		// En QueryRow no podemos retornar error directamente
		// El error se manejará cuando se haga Scan()
		return nil
	}

	return db.QueryRowContext(ctx, query, args...)
}

// EjecutarExec ejecuta un comando SQL (INSERT, UPDATE, DELETE)
func (s *ServicioDB) EjecutarExec(ctx context.Context, query string, usarSecundaria bool, args ...interface{}) (sql.Result, error) {
	db, err := s.ObtenerConexion(usarSecundaria)
	if err != nil {
		return nil, fmt.Errorf("error al obtener conexión: %w", err)
	}

	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar comando: %w", err)
	}

	return result, nil
}

// EjecutarSP ejecuta un Stored Procedure
// Para obtener valores de OUTPUT, usar EjecutarQueryRow con el patrón:
// DECLARE @output INT; EXEC sp_name @param1, @output OUTPUT; SELECT @output
func (s *ServicioDB) EjecutarSP(ctx context.Context, nombreSP string, usarSecundaria bool, args ...interface{}) error {
	db, err := s.ObtenerConexion(usarSecundaria)
	if err != nil {
		return fmt.Errorf("error al obtener conexión: %w", err)
	}

	query := fmt.Sprintf("EXEC %s", nombreSP)
	_, err = db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error al ejecutar SP: %w", err)
	}

	return nil
}

// ObtenerDB retorna la instancia de *sql.DB directamente
// Útil para casos especiales donde se necesita acceso directo
func (s *ServicioDB) ObtenerDB(usarSecundaria bool) (*sql.DB, error) {
	return s.ObtenerConexion(usarSecundaria)
}
