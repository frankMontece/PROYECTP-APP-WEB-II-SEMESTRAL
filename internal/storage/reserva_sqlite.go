package storage

import (
	"condominio-api/internal/models"

	"gorm.io/gorm"
)

// ReservaSQLite implementa ReservaAreaRepository usando GORM sobre SQLite
type ReservaSQLite struct {
	db *gorm.DB
}

// NewReservaSQLite construye el repositorio de reservas inyectándole la conexión ya abierta
func NewReservaSQLite(db *gorm.DB) *ReservaSQLite {
	return &ReservaSQLite{db: db}
}

// Verificación en compilación: si ReservaSQLite no implementa ReservaAreaRepository
var _ ReservaAreaRepository = (*ReservaSQLite)(nil)

// ─── RESERVAS ────────────────────────────────────────────────────────────────

func (s *ReservaSQLite) ListarReservas() []models.ReservaArea {
	var reservas []models.ReservaArea
	s.db.Find(&reservas)
	return reservas
}

func (s *ReservaSQLite) BuscarReservaPorID(id uint) (models.ReservaArea, bool) {
	var r models.ReservaArea
	if err := s.db.First(&r, id).Error; err != nil {
		return models.ReservaArea{}, false
	}
	return r, true
}

func (s *ReservaSQLite) CrearReserva(r models.ReservaArea) models.ReservaArea {
	s.db.Create(&r)
	return r
}

func (s *ReservaSQLite) ActualizarReserva(id uint, datos models.ReservaArea) (models.ReservaArea, bool) {
	var existente models.ReservaArea
	if err := s.db.First(&existente, id).Error; err != nil {
		return models.ReservaArea{}, false
	}
	datos.ID = id
	s.db.Save(&datos)
	return datos, true
}

func (s *ReservaSQLite) BorrarReserva(id uint) bool {
	res := s.db.Delete(&models.ReservaArea{}, id)
	return res.RowsAffected > 0
}
