package service

import (
	"condominio-api/internal/models"
	"condominio-api/internal/storage"
	"strings"
)

// ReservaService contiene la lógica de negocio para Reservas
type ReservaService struct {
	repo storage.ReservaAreaRepository
}

// NewReservaService crea una nueva instancia del servicio de reservas
func NewReservaService(repo storage.ReservaAreaRepository) *ReservaService {
	return &ReservaService{repo: repo}
}

// ListarReservas retorna todas las reservas
func (s *ReservaService) ListarReservas() []models.ReservaArea {
	return s.repo.ListarReservas()
}

// ObtenerReserva retorna una reserva por ID
func (s *ReservaService) ObtenerReserva(id uint) (models.ReservaArea, error) {
	r, ok := s.repo.BuscarReservaPorID(id)
	if !ok {
		return models.ReservaArea{}, ErrNoEncontrado
	}
	return r, nil
}

// CrearReserva valida y crea una reserva
func (s *ReservaService) CrearReserva(r models.ReservaArea) (models.ReservaArea, error) {
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
	return s.repo.CrearReserva(r), nil
}

// ActualizarReserva valida y actualiza una reserva
func (s *ReservaService) ActualizarReserva(id uint, datos models.ReservaArea) (models.ReservaArea, error) {
	r, ok := s.repo.ActualizarReserva(id, datos)
	if !ok {
		return models.ReservaArea{}, ErrNoEncontrado
	}
	return r, nil
}

// BorrarReserva elimina una reserva
func (s *ReservaService) BorrarReserva(id uint) error {
	if !s.repo.BorrarReserva(id) {
		return ErrNoEncontrado
	}
	return nil
}
