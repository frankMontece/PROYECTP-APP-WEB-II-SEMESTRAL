package storage

import (
	"time"

	"gorm.io/gorm"

	"condominio-api/internal/models"
)

// ObligacionSQLite implementa ObligacionRepository usando GORM sobre SQLite.
type ObligacionSQLite struct{ db *gorm.DB }

func NewObligacionSQLite(db *gorm.DB) *ObligacionSQLite {
	return &ObligacionSQLite{db: db}
}

// SembrarVacio inserta datos de ejemplo solo si la tabla está vacía.
// Se llama una vez en main.go justo después de crear el almacén.
func (g *ObligacionSQLite) SembrarVacio() {
	var n int64
	g.db.Model(&models.Obligacion{}).Count(&n)
	if n > 0 {
		return
	}

	ahora := time.Now()
	pagoHace10Dias := ahora.AddDate(0, 0, -10)

	obligaciones := []models.Obligacion{
		{
			ResidenteID:      1,
			Tipo:             "mensual",
			Monto:            85.50,
			Periodo:          ahora.Format("2006-01"),
			Estado:           "pendiente",
			FechaEmision:     ahora.AddDate(0, 0, -15),
			FechaVencimiento: ahora.AddDate(0, 0, 15),
			MoraCalculada:    0,
		},
		{
			ResidenteID:      2,
			Tipo:             "mensual",
			Monto:            85.50,
			Periodo:          ahora.AddDate(0, -1, 0).Format("2006-01"),
			Estado:           "pagada",
			FechaEmision:     ahora.AddDate(0, -1, -15),
			FechaVencimiento: ahora.AddDate(0, -1, 15),
			FechaPago:        &pagoHace10Dias,
			Comprobante:      "TRX-00123",
			MoraCalculada:    0,
		},
		{
			ResidenteID:      3,
			Tipo:             "extraordinaria",
			Monto:            120.00,
			Periodo:          ahora.Format("2006-01"),
			Estado:           "vencida",
			FechaEmision:     ahora.AddDate(0, -2, 0),
			FechaVencimiento: ahora.AddDate(0, -1, 0),
			MoraCalculada:    12.00,
		},
	}
	g.db.Create(&obligaciones)
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

var _ ObligacionRepository = (*ObligacionSQLite)(nil)
