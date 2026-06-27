package storage

import (
	"condominio-api/internal/models"

	"gorm.io/gorm"
)

// MultasGORM implementa MultaRepository usando SQLite + GORM.
type MultasGORM struct{ db *gorm.DB }

func NewMultasGORM(db *gorm.DB) *MultasGORM {
	return &MultasGORM{db: db}
}

func (g *MultasGORM) ListarMultas() []models.Multa {
	var lista []models.Multa
	g.db.Find(&lista)
	return lista
}

func (g *MultasGORM) BuscarMultaPorID(id uint) (models.Multa, bool) {
	var m models.Multa
	if err := g.db.First(&m, id).Error; err != nil {
		return models.Multa{}, false
	}
	return m, true
}

func (g *MultasGORM) CrearMulta(m models.Multa) models.Multa {
	g.db.Create(&m)
	return m
}

func (g *MultasGORM) ActualizarMulta(id uint, datos models.Multa) (models.Multa, bool) {
	var existente models.Multa
	if err := g.db.First(&existente, id).Error; err != nil {
		return models.Multa{}, false
	}
	datos.ID = id
	g.db.Save(&datos)
	return datos, true
}

func (g *MultasGORM) BorrarMulta(id uint) bool {
	return g.db.Delete(&models.Multa{}, id).RowsAffected > 0
}

var _ MultaRepository = (*MultasGORM)(nil)
