package storage

import "condominio-api/internal/models"

// ─── Interfaz de Multa (ISP) ──────────────────────────────────────────────────

type MultaRepository interface {
	ListarMultas() []models.Multa
	BuscarMultaPorID(id uint) (models.Multa, bool)
	CrearMulta(m models.Multa) models.Multa
	ActualizarMulta(id uint, datos models.Multa) (models.Multa, bool)
	BorrarMulta(id uint) bool
}
