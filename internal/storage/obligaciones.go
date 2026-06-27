package storage

import (
	"condominio-api/internal/models"

	"gorm.io/gorm"
)

// ObligacionesGORM implementa ObligacionRepository usando SQLite + GORM.
type ObligacionesGORM struct{ db *gorm.DB }

func NewObligacionesGORM(db *gorm.DB) *ObligacionesGORM {
	return &ObligacionesGORM{db: db}
}

func (g *ObligacionesGORM) ListarObligaciones() []models.Obligacion {
	var lista []models.Obligacion
	g.db.Find(&lista)
	return lista
}

func (g *ObligacionesGORM) BuscarObligacionPorID(id uint) (models.Obligacion, bool) {
	var o models.Obligacion
	if err := g.db.First(&o, id).Error; err != nil {
		return models.Obligacion{}, false
	}
	return o, true
}

func (g *ObligacionesGORM) CrearObligacion(o models.Obligacion) models.Obligacion {
	g.db.Create(&o)
	return o
}

func (g *ObligacionesGORM) ActualizarObligacion(id uint, datos models.Obligacion) (models.Obligacion, bool) {
	var existente models.Obligacion
	if err := g.db.First(&existente, id).Error; err != nil {
		return models.Obligacion{}, false
	}
	datos.ID = id
	g.db.Save(&datos)
	return datos, true
}

func (g *ObligacionesGORM) BorrarObligacion(id uint) bool {
	return g.db.Delete(&models.Obligacion{}, id).RowsAffected > 0
}

var _ ObligacionRepository = (*ObligacionesGORM)(nil)
