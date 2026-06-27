package service

import (
	"strings"

	"condominio-api/internal/models"
	"condominio-api/internal/storage"
)

// MultasService contiene la lógica de negocio de Multa.
type MultasService struct {
	multas storage.MultaRepository
}

func NewMultasService(m storage.MultaRepository) *MultasService {
	return &MultasService{multas: m}
}

func (s *MultasService) ListarMultas() []models.Multa {
	return s.multas.ListarMultas()
}

func (s *MultasService) ObtenerMulta(id uint) (models.Multa, error) {
	m, ok := s.multas.BuscarMultaPorID(id)
	if !ok {
		return models.Multa{}, ErrNoEncontrado
	}
	return m, nil
}

func (s *MultasService) CrearMulta(m models.Multa) (models.Multa, error) {
	if m.ResidenteID == 0 {
		return models.Multa{}, ErrResidenteIDInvalido
	}
	if strings.TrimSpace(m.Motivo) == "" {
		return models.Multa{}, ErrMotivoVacio
	}
	if m.Monto <= 0 {
		return models.Multa{}, ErrMontoInvalido
	}
	m.Estado = "pendiente"
	return s.multas.CrearMulta(m), nil
}

func (s *MultasService) ActualizarMulta(id uint, datos models.Multa) (models.Multa, error) {
	if datos.ResidenteID == 0 {
		return models.Multa{}, ErrResidenteIDInvalido
	}
	if strings.TrimSpace(datos.Motivo) == "" {
		return models.Multa{}, ErrMotivoVacio
	}
	if datos.Monto <= 0 {
		return models.Multa{}, ErrMontoInvalido
	}
	m, ok := s.multas.ActualizarMulta(id, datos)
	if !ok {
		return models.Multa{}, ErrNoEncontrado
	}
	return m, nil
}

func (s *MultasService) BorrarMulta(id uint) error {
	if !s.multas.BorrarMulta(id) {
		return ErrNoEncontrado
	}
	return nil
}
