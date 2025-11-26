package atenciones

import (
	"context"
	"database/sql"
	"fmt"

	"backend/internal/shared/database"
)

type InfoFacturacionAtencion struct {
	IdPaciente             int
	IdServicio             int
	IdFuenteFinanciamiento int
	IdTipoFinanciamiento   int
	IdEstadoFacturacion    int
	TieneHemoglobina       bool
}

type DatosPaciente struct {
	EdadPaciente       *int
	NroHistoriaClinica *float64
	NombreMedico       *string
	IdServicio         *int
	NombreServicio     *string
}

type AtencionesServicio struct {
	db *database.ServicioDB
}

func NuevoServicio(db *database.ServicioDB) *AtencionesServicio {
	return &AtencionesServicio{db: db}
}

func (s *AtencionesServicio) ObtenerInfoFacturacionAtencion(ctx context.Context, idAtencion int) (*InfoFacturacionAtencion, error) {
	row := s.db.EjecutarQueryRow(ctx, QueryObtenerInfoFacturacionAtencion, false, sql.Named("idAtencion", idAtencion))
	if row == nil {
		return nil, fmt.Errorf("error al obtener conexión a la base de datos")
	}

	var info InfoFacturacionAtencion
	err := row.Scan(
		&info.IdPaciente,
		&info.IdServicio,
		&info.IdFuenteFinanciamiento,
		&info.IdTipoFinanciamiento,
		&info.IdEstadoFacturacion,
		&info.TieneHemoglobina,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no se encontró información del N° Cuenta %d", idAtencion)
	}
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (s *AtencionesServicio) ObtenerDatosPaciente(ctx context.Context, idAtencion int) (*DatosPaciente, error) {
	row := s.db.EjecutarQueryRow(ctx, QueryObtenerDatosPaciente, false, sql.Named("idAtencion", idAtencion))
	if row == nil {
		return nil, fmt.Errorf("error al obtener conexión a la base de datos")
	}

	var datos DatosPaciente
	err := row.Scan(
		&datos.EdadPaciente,
		&datos.NroHistoriaClinica,
		&datos.NombreMedico,
		&datos.IdServicio,
		&datos.NombreServicio,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("paciente no encontrado en el N° Cuenta: %d", idAtencion)
	}
	if err != nil {
		return nil, err
	}

	return &datos, nil
}
