package storage

import (
	"condominio-api/internal/models"

	"gorm.io/gorm"
)

// ObligacionSQLite implementa ObligacionRepository usando GORM sobre SQLite.
type ObligacionSQLite struct{ db *gorm.DB }

func NewObligacionSQLite(db *gorm.DB) *ObligacionSQLite {
	return &ObligacionSQLite{db: db}
}

func (g *ObligacionSQLite) ListarObligaciones() []models.Obligacion {
	var lista []models.Obligacion
	g.db.Find(&lista)
	return lista
}

func (g *ObligacionSQLite) BuscarObligacionPorID(id uint) (models.Obligacion, bool) {
	var o models.Obligacion
	if err := g.db.First(&o, id).Error; err != nil {
		return models.Obligacion{}, false
	}
	return o, true
}

func (g *ObligacionSQLite) CrearObligacion(o models.Obligacion) models.Obligacion {
	g.db.Create(&o)
	return o
}

func (g *ObligacionSQLite) ActualizarObligacion(id uint, datos models.Obligacion) (models.Obligacion, bool) {
	var existente models.Obligacion
	if err := g.db.First(&existente, id).Error; err != nil {
		return models.Obligacion{}, false
	}
	datos.ID = id
	g.db.Save(&datos)
	return datos, true
}

func (g *ObligacionSQLite) BorrarObligacion(id uint) bool {
	return g.db.Delete(&models.Obligacion{}, id).RowsAffected > 0
}

// Verificación en compilación: si ObligacionSQLite no implementa
// todos los métodos de ObligacionRepository, el compilador falla aquí.
var _ ObligacionRepository = (*ObligacionSQLite)(nil)
