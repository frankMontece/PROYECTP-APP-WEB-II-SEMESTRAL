package storage

import (
	"condominio-api/internal/models"
	"time"

	"gorm.io/gorm"
)

// VisitaSQLite implementa VisitasRepository usando GORM sobre SQLite
type VisitaSQLite struct {
	db *gorm.DB
}

// NewVisitaSQLite construye el repositorio de visitas inyectándole la conexión ya abierta
func NewVisitaSQLite(db *gorm.DB) *VisitaSQLite {
	return &VisitaSQLite{db: db}
}

// VISITAS

func (s *VisitaSQLite) ListarVisitas() []models.VisitaVehiculo {
	var visitas []models.VisitaVehiculo
	s.db.Find(&visitas)
	return visitas
}

func (s *VisitaSQLite) BuscarVisitaPorID(id uint) (models.VisitaVehiculo, bool) {
	var vis models.VisitaVehiculo
	if err := s.db.First(&vis, id).Error; err != nil {
		return models.VisitaVehiculo{}, false
	}
	return vis, true
}

func (s *VisitaSQLite) CrearVisita(vis models.VisitaVehiculo) models.VisitaVehiculo {
	s.db.Create(&vis)
	return vis
}

func (s *VisitaSQLite) ActualizarVisita(id uint, datos models.VisitaVehiculo) (models.VisitaVehiculo, bool) {
	var existente models.VisitaVehiculo
	if err := s.db.First(&existente, id).Error; err != nil {
		return models.VisitaVehiculo{}, false
	}
	datos.ID = id
	s.db.Save(&datos)
	return datos, true
}

// RegistrarSalidaVisita marca la hora de salida y expira el QR.
// Usamos s.db.Save para que el cambio persista en disco:
// el próximo reinicio del servidor conservará el registro de salida.
func (s *VisitaSQLite) RegistrarSalidaVisita(id uint) (models.VisitaVehiculo, bool) {
	var vis models.VisitaVehiculo
	if err := s.db.First(&vis, id).Error; err != nil {
		return models.VisitaVehiculo{}, false
	}
	ahora := time.Now()
	vis.HoraSalida = &ahora
	vis.EstadoQR = "expirado"
	s.db.Save(&vis)
	return vis, true
}

func (s *VisitaSQLite) BorrarVisita(id uint) bool {
	res := s.db.Delete(&models.VisitaVehiculo{}, id)
	return res.RowsAffected > 0
}
