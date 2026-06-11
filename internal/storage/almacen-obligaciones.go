package storage

import "condominio-api/internal/models"

// AlmacenObligaciones define el contrato que debe cumplir
// cualquier implementación de almacenamiento del módulo de alícuotas.
type AlmacenObligaciones interface {
	// Obligaciones
	ListarObligaciones() ([]models.Obligacion, error)
	BuscarObligacionPorID(id uint) (*models.Obligacion, error)
	CrearObligacion(o models.Obligacion) (*models.Obligacion, error)
	ActualizarObligacion(id uint, datos models.Obligacion) (*models.Obligacion, error)
	EliminarObligacion(id uint) error

	// Multas
	ListarMultas() ([]models.Multa, error)
	BuscarMultaPorID(id uint) (*models.Multa, error)
	CrearMulta(m models.Multa) (*models.Multa, error)
	ActualizarMulta(id uint, datos models.Multa) (*models.Multa, error)
	EliminarMulta(id uint) error
}

// Verificación en tiempo de compilación: MemoriaObligaciones debe implementar AlmacenObligaciones.
// Si falta algún método, Go lanza error aquí antes de correr.
var _ AlmacenObligaciones = (*MemoriaObligaciones)(nil)
