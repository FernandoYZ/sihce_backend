package auditoria

import (
	"context"
	"database/sql"

	"backend/internal/shared/database"
)

// Constantes de acciones de auditoría
const (
	AccionAgregar   = "A"
	AccionModificar = "M"
	AccionEliminar  = "E"
)

type AuditoriaServicio struct {
	db *database.ServicioDB
}

func NuevoServicio(db *database.ServicioDB) *AuditoriaServicio {
	return &AuditoriaServicio{db: db}
}

func (s *AuditoriaServicio) ObtenerNombreEmpleado(ctx context.Context, idUsuario int) (string, error) {
	row := s.db.EjecutarQueryRow(ctx, QueryObtenerNombreCompleto, false, sql.Named("idUsuario", idUsuario))
	if row == nil {
		return "API", nil
	}

	var usuario sql.NullString
	err := row.Scan(&usuario)

	if err == sql.ErrNoRows || !usuario.Valid {
		return "API", nil
	}
	if err != nil {
		return "API", err
	}

	return usuario.String, nil
}

func (s *AuditoriaServicio) RegistrarAuditoria(
	ctx context.Context,
	idEmpleado int,
	accion string,
	idRegistro int,
	tabla string,
	idListItem int,
	nombrePC string,
	observaciones string,
) error {
	return s.db.EjecutarSP(
		ctx,
		"AuditoriaAgregarV @IdEmpleado, @Accion, @IdRegistro, @Tabla, @idListItem, @nombrePC, @observaciones",
		false,
		sql.Named("IdEmpleado", idEmpleado),
		sql.Named("Accion", accion),
		sql.Named("IdRegistro", idRegistro),
		sql.Named("Tabla", tabla),
		sql.Named("idListItem", idListItem),
		sql.Named("nombrePC", nombrePC),
		sql.Named("observaciones", observaciones),
	)
}
