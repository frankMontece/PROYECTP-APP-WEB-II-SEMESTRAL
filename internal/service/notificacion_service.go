package service

import (
	"condominio-api/internal/models"
	"condominio-api/internal/storage"
	"strings"
)

// NotificacionService contiene la lógica de negocio para Notificaciones
type NotificacionService struct {
	repo storage.NotificacionRepository
}

// NewNotificacionService crea una nueva instancia del servicio de notificaciones
func NewNotificacionService(repo storage.NotificacionRepository) *NotificacionService {
	return &NotificacionService{repo: repo}
}

// ListarNotificaciones retorna todas las notificaciones
func (s *NotificacionService) ListarNotificaciones() []models.Notificacion {
	return s.repo.ListarNotificaciones()
}

// ObtenerNotificacion retorna una notificación por ID
func (s *NotificacionService) ObtenerNotificacion(id uint) (models.Notificacion, error) {
	n, ok := s.repo.BuscarNotificacionPorID(id)
	if !ok {
		return models.Notificacion{}, ErrNoEncontrado
	}
	return n, nil
}

// CrearNotificacion valida y crea una notificación
func (s *NotificacionService) CrearNotificacion(n models.Notificacion) (models.Notificacion, error) {
	if n.ResidenteID == 0 {
		return models.Notificacion{}, ErrNoEncontrado
	}
	if strings.TrimSpace(n.Tipo) == "" || strings.TrimSpace(n.Mensaje) == "" {
		return models.Notificacion{}, ErrNombreVacio
	}
	return s.repo.CrearNotificacion(n), nil
}

// MarcarLeida marca una notificación como leída
func (s *NotificacionService) MarcarLeida(id uint) (models.Notificacion, error) {
	n, ok := s.repo.MarcarComoLeida(id)
	if !ok {
		return models.Notificacion{}, ErrNoEncontrado
	}
	return n, nil
}

// BorrarNotificacion elimina una notificación
func (s *NotificacionService) BorrarNotificacion(id uint) error {
	if !s.repo.BorrarNotificacion(id) {
		return ErrNoEncontrado
	}
	return nil
}
