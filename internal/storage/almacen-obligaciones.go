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

type MultaRepository interface {
	ListarMultas() []models.Multa
	BuscarMultaPorID(id uint) (models.Multa, bool)
	CrearMulta(m models.Multa) models.Multa
	ActualizarMulta(id uint, datos models.Multa) (models.Multa, bool)
	BorrarMulta(id uint) bool
}

type AlmacenObligaciones interface {
	ObligacionRepository
	MultaRepository
}

var _ ObligacionRepository = (*ObligacionSQLite)(nil)

var _ MultaRepository = (*MultaSQLite)(nil)
