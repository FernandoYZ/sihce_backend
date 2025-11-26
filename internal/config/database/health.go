package database

import (
	"context"
	"fmt"
)

// VerificarSalud verifica el estado de las conexiones activas
func (g *GestorDB) VerificarSalud(ctx context.Context) error {
	g.mu.RLock()
	principal := g.principal
	secundaria := g.secundaria
	g.mu.RUnlock()

	// Verificación BD Principal (Crítica)
	if principal == nil {
		// En un health check, si la conexión principal no está activa (porque no se ha usado aún),
		// debería considerarse un fallo si se espera que esté siempre disponible tras el arranque.
		// Si la inicialización es puramente lazy, este error podría ser opcional.
		// Dado que existe `InicializarPrincipal`, asumimos que es crítico.
		return fmt.Errorf("la base de datos principal no está inicializada")
	}

	if err := principal.PingContext(ctx); err != nil {
		return fmt.Errorf("error en health check de base de datos principal: %w", err)
	}

	// Verificación BD Secundaria (Opcional - Lazy)
	if secundaria != nil {
		if err := secundaria.PingContext(ctx); err != nil {
			return fmt.Errorf("error en health check de base de datos secundaria: %w", err)
		}
	}

	return nil
}
