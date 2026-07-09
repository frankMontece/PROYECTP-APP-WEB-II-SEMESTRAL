package storage

import (
	"fmt"
	"time"

	"condominio-api/internal/models"

	"gorm.io/gorm"
)

func SembrarObligacion(db *gorm.DB) {
	var n int64
	db.Model(&models.Obligacion{}).Count(&n)
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
	db.Create(&obligaciones)

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
	db.Create(&multas)

	fmt.Printf("✅ Datos de ejemplo sembrados en Obligaciones: %d obligaciones, %d multas\n",
		len(obligaciones), len(multas))
}
