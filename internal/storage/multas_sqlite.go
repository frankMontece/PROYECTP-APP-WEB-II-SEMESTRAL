package storage

import (
	"time"

	"condominio-api/internal/models"

	"gorm.io/gorm"
)

// MultaSQLite implementa MultaRepository usando GORM sobre SQLite.
type MultaSQLite struct{ db *gorm.DB }

func NewMultaSQLite(db *gorm.DB) *MultaSQLite {
	return &MultaSQLite{db: db}
}

// SembrarVacio inserta datos de ejemplo solo si la tabla está vacía.
// Se llama una vez en main.go justo después de crear el almacén.
// Si ya hay datos, no hace nada: es seguro llamarlo en cada arranque.
func (g *MultaSQLite) SembrarVacio() {
	var n int64
	g.db.Model(&models.Multa{}).Count(&n)
	if n > 0 {
		return
	}

	ahora := time.Now()
	pagoHace5Dias := ahora.AddDate(0, 0, -5)

	multas := []models.Multa{
		{
			ResidenteID:  1,
			Motivo:       "Ruido excesivo",
			Monto:        30.00,
			Estado:       "pendiente",
			FechaEmision: ahora.AddDate(0, 0, -3),
		},
		{
			ResidenteID:  2,
			Motivo:       "Uso indebido de áreas comunes",
			Monto:        20.00,
			Estado:       "pagada",
			FechaEmision: ahora.AddDate(0, 0, -8),
			FechaPago:    &pagoHace5Dias,
		},
		{
			ResidenteID:  3,
			Motivo:       "Mascota sin correa",
			Monto:        15.00,
			Estado:       "apelada",
			FechaEmision: ahora.AddDate(0, 0, -1),
		},
	}
	g.db.Create(&multas)
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
