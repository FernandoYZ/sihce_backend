package auditoria

const (
	QueryObtenerNombreCompleto = `
  SELECT
    LEFT(
      UPPER(LTRIM(RTRIM(ApellidoPaterno + ' ' + ISNULL(ApellidoMaterno, '') + ' ' + Nombres))),
      30
    ) AS Usuario
  FROM Empleados
  WHERE IdEmpleado = @idUsuario`
)
