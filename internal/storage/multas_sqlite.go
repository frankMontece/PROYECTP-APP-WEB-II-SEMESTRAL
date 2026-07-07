package storage

import (
	"condominio-api/internal/models"

	"gorm.io/gorm"
)

// MultaSQLite implementa MultaRepository usando GORM sobre SQLite.
type MultaSQLite struct{ db *gorm.DB }

func NewMultaSQLite(db *gorm.DB) *MultaSQLite {
	return &MultaSQLite{db: db}
}

func (g *MultaSQLite) ListarMultas() []models.Multa {
	var lista []models.Multa
	g.db.Find(&lista)
	return lista
}

func (g *MultaSQLite) BuscarMultaPorID(id uint) (models.Multa, bool) {
	var m models.Multa
	if err := g.db.First(&m, id).Error; err != nil {
		return models.Multa{}, false
	}
	return m, true
}

func (g *MultaSQLite) CrearMulta(m models.Multa) models.Multa {
	g.db.Create(&m)
	return m
}

func (g *MultaSQLite) ActualizarMulta(id uint, datos models.Multa) (models.Multa, bool) {
	var existente models.Multa
	if err := g.db.First(&existente, id).Error; err != nil {
		return models.Multa{}, false
	}
	datos.ID = id
	g.db.Save(&datos)
	return datos, true
}

func (g *MultaSQLite) BorrarMulta(id uint) bool {
	return g.db.Delete(&models.Multa{}, id).RowsAffected > 0
}

// Verificación en compilación: si MultaSQLite no implementa
// todos los métodos de MultaRepository, el compilador falla aquí.
var _ MultaRepository = (*MultaSQLite)(nil)
