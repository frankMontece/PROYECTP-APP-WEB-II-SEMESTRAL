package service

import (
	"condominio-api/internal/models"
	"condominio-api/internal/storage"
	"strings"
)

// SocialService contiene la lógica de negocio del módulo B
type SocialService struct {
	areas    storage.AreaSocialRepository
	reservas storage.ReservaAreaRepository
	notifs   storage.NotificacionRepository
}

// NewSocialService crea una nueva instancia del servicio social
func NewSocialService(
	a storage.AreaSocialRepository,
	r storage.ReservaAreaRepository,
	n storage.NotificacionRepository,
) *SocialService {
	return &SocialService{areas: a, reservas: r, notifs: n}
}

// ─── ÁREAS SOCIALES ──────────────────────────────────────────────────────────

func (s *SocialService) ListarAreas() []models.AreaSocial {
	return s.areas.ListarAreas()
}

func (s *SocialService) ObtenerArea(id uint) (models.AreaSocial, error) {
	a, ok := s.areas.BuscarAreaPorID(id)
	if !ok {
		return models.AreaSocial{}, ErrNoEncontrado
	}
	return a, nil
}

func (s *SocialService) CrearArea(a models.AreaSocial) (models.AreaSocial, error) {
	if strings.TrimSpace(a.Nombre) == "" {
		return models.AreaSocial{}, ErrNombreVacio
	}
	if a.Capacidad <= 0 {
		return models.AreaSocial{}, ErrCapacidadInvalida
	}
	a.Activo = true
	return s.areas.CrearArea(a), nil
}

func (s *SocialService) ActualizarArea(id uint, datos models.AreaSocial) (models.AreaSocial, error) {
	if strings.TrimSpace(datos.Nombre) == "" {
		return models.AreaSocial{}, ErrNombreVacio
	}
	a, ok := s.areas.ActualizarArea(id, datos)
	if !ok {
		return models.AreaSocial{}, ErrNoEncontrado
	}
	return a, nil
}

func (s *SocialService) BorrarArea(id uint) error {
	if !s.areas.BorrarArea(id) {
		return ErrNoEncontrado
	}
	return nil
}

// ─── RESERVAS ────────────────────────────────────────────────────────────────

func (s *SocialService) ListarReservas() []models.ReservaArea {
	return s.reservas.ListarReservas()
}

func (s *SocialService) ObtenerReserva(id uint) (models.ReservaArea, error) {
	r, ok := s.reservas.BuscarReservaPorID(id)
	if !ok {
		return models.ReservaArea{}, ErrNoEncontrado
	}
	return r, nil
}

func (s *SocialService) CrearReserva(r models.ReservaArea) (models.ReservaArea, error) {
	if r.AreaID == 0 || r.ResidenteID == 0 {
		return models.ReservaArea{}, ErrNoEncontrado
	}
	if strings.TrimSpace(r.Proposito) == "" {
		return models.ReservaArea{}, ErrNombreVacio
	}
	if r.NumeroPersonas <= 0 {
		return models.ReservaArea{}, ErrNumeroPersonasInvalido
	}
	if r.FechaInicio.IsZero() || r.FechaFin.IsZero() || !r.FechaFin.After(r.FechaInicio) {
		return models.ReservaArea{}, ErrFechasInvalidas
	}
	r.Estado = "pendiente"
	return s.reservas.CrearReserva(r), nil
}

func (s *SocialService) ActualizarReserva(id uint, datos models.ReservaArea) (models.ReservaArea, error) {
	r, ok := s.reservas.ActualizarReserva(id, datos)
	if !ok {
		return models.ReservaArea{}, ErrNoEncontrado
	}
	return r, nil
}

func (s *SocialService) BorrarReserva(id uint) error {
	if !s.reservas.BorrarReserva(id) {
		return ErrNoEncontrado
	}
	return nil
}

// ─── NOTIFICACIONES ──────────────────────────────────────────────────────────

func (s *SocialService) ListarNotificaciones() []models.Notificacion {
	return s.notifs.ListarNotificaciones()
}

func (s *SocialService) ObtenerNotificacion(id uint) (models.Notificacion, error) {
	n, ok := s.notifs.BuscarNotificacionPorID(id)
	if !ok {
		return models.Notificacion{}, ErrNoEncontrado
	}
	return n, nil
}

func (s *SocialService) CrearNotificacion(n models.Notificacion) (models.Notificacion, error) {
	if n.ResidenteID == 0 {
		return models.Notificacion{}, ErrNoEncontrado
	}
	if strings.TrimSpace(n.Tipo) == "" || strings.TrimSpace(n.Mensaje) == "" {
		return models.Notificacion{}, ErrNombreVacio
	}
	return s.notifs.CrearNotificacion(n), nil
}

func (s *SocialService) MarcarLeida(id uint) (models.Notificacion, error) {
	n, ok := s.notifs.MarcarComoLeida(id)
	if !ok {
		return models.Notificacion{}, ErrNoEncontrado
	}
	return n, nil
}

func (s *SocialService) BorrarNotificacion(id uint) error {
	if !s.notifs.BorrarNotificacion(id) {
		return ErrNoEncontrado
	}
	return nil
}
