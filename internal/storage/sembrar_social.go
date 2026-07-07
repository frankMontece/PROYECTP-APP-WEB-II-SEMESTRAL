package storage

import (
	"fmt"
	"time"

	"condominio-api/internal/models"

	"gorm.io/gorm"
)

func SembrarSocial(db *gorm.DB) {
	var n int64
	db.Model(&models.AreaSocial{}).Count(&n)
	if n > 0 {
		return
	}

	ahora := time.Now()

	// ── Áreas Sociales de ejemplo ──────────────────────────────────────────
	areas := []models.AreaSocial{
		{Nombre: "Salón de Eventos", Descripcion: "Amplio salón con capacidad para eventos", Capacidad: 80, Activo: true},
		{Nombre: "Cancha de Fútbol", Descripcion: "Cancha sintética iluminada", Capacidad: 30, Activo: true},
		{Nombre: "Piscina", Descripcion: "Piscina olímpica con área de descanso", Capacidad: 50, Activo: true},
		{Nombre: "Gimnasio", Descripcion: "Equipamiento completo de pesas y cardio", Capacidad: 20, Activo: true},
	}
	db.Create(&areas)

	// ── Reservas de ejemplo ────────────────────────────────────────────────
	reservas := []models.ReservaArea{
		{
			AreaID:         1,
			ResidenteID:    1,
			FechaInicio:    ahora.Add(24 * time.Hour),
			FechaFin:       ahora.Add(26 * time.Hour),
			Proposito:      "Cumpleaños de 15 años",
			NumeroPersonas: 60,
			Estado:         "aprobada",
			FechaCreacion:  ahora.Add(-2 * time.Hour),
		},
		{
			AreaID:         2,
			ResidenteID:    2,
			FechaInicio:    ahora.Add(48 * time.Hour),
			FechaFin:       ahora.Add(50 * time.Hour),
			Proposito:      "Partido amistoso",
			NumeroPersonas: 22,
			Estado:         "pendiente",
			FechaCreacion:  ahora.Add(-1 * time.Hour),
		},
		{
			AreaID:         3,
			ResidenteID:    3,
			FechaInicio:    ahora.Add(72 * time.Hour),
			FechaFin:       ahora.Add(74 * time.Hour),
			Proposito:      "Clase de natación",
			NumeroPersonas: 15,
			Estado:         "pendiente",
			FechaCreacion:  ahora,
		},
	}
	db.Create(&reservas)

	// ── Notificaciones de ejemplo ──────────────────────────────────────────
	notificaciones := []models.Notificacion{
		{
			ResidenteID:   1,
			Tipo:          "reserva",
			Mensaje:       "Su reserva para el Salón de Eventos ha sido aprobada",
			Leida:         false,
			FechaCreacion: ahora.Add(-1 * time.Hour),
		},
		{
			ResidenteID:   2,
			Tipo:          "aviso",
			Mensaje:       "Recordatorio: su reserva para la Cancha de Fútbol está pendiente de aprobación",
			Leida:         false,
			FechaCreacion: ahora.Add(-30 * time.Minute),
		},
		{
			ResidenteID:   3,
			Tipo:          "pago",
			Mensaje:       "Su pago de mantenimiento de junio está vencido",
			Leida:         true,
			FechaCreacion: ahora.Add(-5 * time.Hour),
		},
	}
	db.Create(&notificaciones)

	fmt.Printf("✅ Datos de ejemplo sembrados en Social: %d áreas, %d reservas, %d notificaciones\n",
		len(areas), len(reservas), len(notificaciones))
}
