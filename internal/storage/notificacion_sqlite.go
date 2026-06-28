package storage

import (
	"condominio-api/internal/models"

	"gorm.io/gorm"
)

// NotificacionSQLite implementa NotificacionRepository usando GORM sobre SQLite
type NotificacionSQLite struct {
	db *gorm.DB
}

// NewNotificacionSQLite construye el repositorio de notificaciones inyectándole la conexión ya abierta
func NewNotificacionSQLite(db *gorm.DB) *NotificacionSQLite {
	return &NotificacionSQLite{db: db}
}

// Verificación en compilación: si NotificacionSQLite no implementa NotificacionRepository
var _ NotificacionRepository = (*NotificacionSQLite)(nil)

// ─── NOTIFICACIONES ──────────────────────────────────────────────────────────

func (s *NotificacionSQLite) ListarNotificaciones() []models.Notificacion {
	var notifs []models.Notificacion
	s.db.Find(&notifs)
	return notifs
}

func (s *NotificacionSQLite) BuscarNotificacionPorID(id uint) (models.Notificacion, bool) {
	var n models.Notificacion
	if err := s.db.First(&n, id).Error; err != nil {
		return models.Notificacion{}, false
	}
	return n, true
}

func (s *NotificacionSQLite) CrearNotificacion(n models.Notificacion) models.Notificacion {
	s.db.Create(&n)
	return n
}

func (s *NotificacionSQLite) MarcarComoLeida(id uint) (models.Notificacion, bool) {
	var n models.Notificacion
	if err := s.db.First(&n, id).Error; err != nil {
		return models.Notificacion{}, false
	}
	n.Leida = true
	s.db.Save(&n)
	return n, true
}

func (s *NotificacionSQLite) BorrarNotificacion(id uint) bool {
	res := s.db.Delete(&models.Notificacion{}, id)
	return res.RowsAffected > 0
}
