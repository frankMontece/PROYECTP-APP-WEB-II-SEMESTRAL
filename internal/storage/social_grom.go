package storage

import (
	"condominio-api/internal/models"

	"gorm.io/gorm"
)

// SocialGORM implementa los repositorios usando GORM + SQLite
type SocialGORM struct {
	db *gorm.DB
}

// NewSocialGORM crea una nueva instancia del repositorio GORM
func NewSocialGORM(db *gorm.DB) *SocialGORM {
	return &SocialGORM{db: db}
}

// ─── ÁREAS SOCIALES ──────────────────────────────────────────────────────────

func (s *SocialGORM) ListarAreas() []models.AreaSocial {
	var areas []models.AreaSocial
	s.db.Find(&areas)
	return areas
}

func (s *SocialGORM) BuscarAreaPorID(id uint) (models.AreaSocial, bool) {
	var a models.AreaSocial
	if err := s.db.First(&a, id).Error; err != nil {
		return models.AreaSocial{}, false
	}
	return a, true
}

func (s *SocialGORM) CrearArea(a models.AreaSocial) models.AreaSocial {
	s.db.Create(&a)
	return a
}

func (s *SocialGORM) ActualizarArea(id uint, datos models.AreaSocial) (models.AreaSocial, bool) {
	var existente models.AreaSocial
	if err := s.db.First(&existente, id).Error; err != nil {
		return models.AreaSocial{}, false
	}
	datos.ID = id
	s.db.Save(&datos)
	return datos, true
}

func (s *SocialGORM) BorrarArea(id uint) bool {
	result := s.db.Delete(&models.AreaSocial{}, id)
	return result.RowsAffected > 0
}

// ─── RESERVAS ────────────────────────────────────────────────────────────────

func (s *SocialGORM) ListarReservas() []models.ReservaArea {
	var reservas []models.ReservaArea
	s.db.Find(&reservas)
	return reservas
}

func (s *SocialGORM) BuscarReservaPorID(id uint) (models.ReservaArea, bool) {
	var r models.ReservaArea
	if err := s.db.First(&r, id).Error; err != nil {
		return models.ReservaArea{}, false
	}
	return r, true
}

func (s *SocialGORM) CrearReserva(r models.ReservaArea) models.ReservaArea {
	s.db.Create(&r)
	return r
}

func (s *SocialGORM) ActualizarReserva(id uint, datos models.ReservaArea) (models.ReservaArea, bool) {
	var existente models.ReservaArea
	if err := s.db.First(&existente, id).Error; err != nil {
		return models.ReservaArea{}, false
	}
	datos.ID = id
	s.db.Save(&datos)
	return datos, true
}

func (s *SocialGORM) BorrarReserva(id uint) bool {
	result := s.db.Delete(&models.ReservaArea{}, id)
	return result.RowsAffected > 0
}

// ─── NOTIFICACIONES ──────────────────────────────────────────────────────────

func (s *SocialGORM) ListarNotificaciones() []models.Notificacion {
	var notifs []models.Notificacion
	s.db.Find(&notifs)
	return notifs
}

func (s *SocialGORM) BuscarNotificacionPorID(id uint) (models.Notificacion, bool) {
	var n models.Notificacion
	if err := s.db.First(&n, id).Error; err != nil {
		return models.Notificacion{}, false
	}
	return n, true
}

func (s *SocialGORM) CrearNotificacion(n models.Notificacion) models.Notificacion {
	s.db.Create(&n)
	return n
}

func (s *SocialGORM) MarcarComoLeida(id uint) (models.Notificacion, bool) {
	var n models.Notificacion
	if err := s.db.First(&n, id).Error; err != nil {
		return models.Notificacion{}, false
	}
	n.Leida = true
	s.db.Save(&n)
	return n, true
}

func (s *SocialGORM) BorrarNotificacion(id uint) bool {
	result := s.db.Delete(&models.Notificacion{}, id)
	return result.RowsAffected > 0
}

// ─── VERIFICACIONES EN TIEMPO DE COMPILACIÓN ──────────────────────────────

var _ AreaSocialRepository = (*SocialGORM)(nil)
var _ ReservaAreaRepository = (*SocialGORM)(nil)
var _ NotificacionRepository = (*SocialGORM)(nil)
