package storage

import "condominio-api/internal/models"

// ─── Interfaz de Obligacion (ISP) ────────────────────────────────────────────

type ObligacionRepository interface {
	ListarObligaciones() []models.Obligacion
	BuscarObligacionPorID(id uint) (models.Obligacion, bool)
	CrearObligacion(o models.Obligacion) models.Obligacion
	ActualizarObligacion(id uint, datos models.Obligacion) (models.Obligacion, bool)
	BorrarObligacion(id uint) bool
}
