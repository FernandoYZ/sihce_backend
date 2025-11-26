package atenciones

const (
	QueryObtenerInfoFacturacionAtencion = `
  SELECT TOP 1
    c.IdPaciente,
    c.IdServicio,
    a.idFuenteFinanciamiento,
    fas.idTipoFinanciamiento,
    fas.IdEstadoFacturacion,
    CASE WHEN EXISTS (
      SELECT 1
      FROM FacturacionServicioDespacho fsd
      INNER JOIN FactOrdenServicio fas2 ON fsd.idOrden = fas2.IdOrden
      INNER JOIN Atenciones a2 ON fas2.IdCuentaAtencion = a2.IdCuentaAtencion
      INNER JOIN Citas c2 ON a2.IdAtencion = c2.IdAtencion
      WHERE c2.IdAtencion = @idAtencion AND fsd.IdProducto = 3588
    ) THEN CAST(1 AS BIT) ELSE CAST(0 AS BIT) END AS TieneHemoglobina
  FROM Citas c
  INNER JOIN Atenciones a ON c.IdAtencion = a.IdAtencion
  INNER JOIN FactOrdenServicio fas ON a.IdCuentaAtencion = fas.IdCuentaAtencion
  WHERE c.IdAtencion = @idAtencion`

	QueryObtenerDatosPaciente = `
  SELECT
    a.Edad as edadPaciente,
    pa.NroHistoriaClinica,
    e.ApellidoPaterno + ' ' + isnull(e.ApellidoMaterno, '') + ' ' + e.Nombres AS NombreMedico,
    s.IdServicio,
    s.Nombre AS nombreServicio
  FROM Atenciones a
  INNER JOIN Citas c ON a.IdAtencion = c.IdAtencion
  INNER JOIN Pacientes pa ON c.IdPaciente = pa.IdPaciente
  INNER JOIN ProgramacionMedica p ON c.IdProgramacion = p.IdProgramacion
  INNER JOIN Servicios s ON p.IdServicio = s.IdServicio
  INNER JOIN Medicos m ON p.IdMedico = m.IdMedico
  INNER JOIN Empleados e ON m.IdEmpleado = e.IdEmpleado
  WHERE a.IdAtencion = @idAtencion`
)
