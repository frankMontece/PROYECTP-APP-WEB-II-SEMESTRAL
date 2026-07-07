package service

import (
	"strings"

	"condominio-api/internal/models"
	"condominio-api/internal/storage"
)

// ObligacionesService contiene la lógica de negocio de Obligacion.
type ObligacionesService struct {
	obligaciones storage.ObligacionRepository
}

func NewObligacionesService(o storage.ObligacionRepository) *ObligacionesService {
	return &ObligacionesService{obligaciones: o}
}

func (s *ObligacionesService) ListarObligaciones() []models.Obligacion {
	return s.obligaciones.ListarObligaciones()
}

func (s *ObligacionesService) ObtenerObligacion(id uint) (models.Obligacion, error) {
	o, ok := s.obligaciones.BuscarObligacionPorID(id)
	if !ok {
		return models.Obligacion{}, ErrNoEncontrado
	}
	return o, nil
}

func (s *ObligacionesService) CrearObligacion(o models.Obligacion) (models.Obligacion, error) {
	if o.ResidenteID == 0 {
		return models.Obligacion{}, ErrResidenteIDInvalido
	}
	if o.Monto <= 0 {
		return models.Obligacion{}, ErrMontoInvalido
	}
	if strings.TrimSpace(o.Periodo) == "" {
		return models.Obligacion{}, ErrPeriodoVacio
	}
	if o.Tipo != "mensual" && o.Tipo != "extraordinaria" {
		return models.Obligacion{}, ErrTipoObligacion
	}
	o.Estado = "pendiente"
	return s.obligaciones.CrearObligacion(o), nil
}

func (s *ObligacionesService) ActualizarObligacion(id uint, datos models.Obligacion) (models.Obligacion, error) {
	if datos.ResidenteID == 0 {
		return models.Obligacion{}, ErrResidenteIDInvalido
	}
	if datos.Monto <= 0 {
		return models.Obligacion{}, ErrMontoInvalido
	}
	if strings.TrimSpace(datos.Periodo) == "" {
		return models.Obligacion{}, ErrPeriodoVacio
	}
	if datos.Tipo != "mensual" && datos.Tipo != "extraordinaria" {
		return models.Obligacion{}, ErrTipoObligacion
	}
	o, ok := s.obligaciones.ActualizarObligacion(id, datos)
	if !ok {
		return models.Obligacion{}, ErrNoEncontrado
	}
	return o, nil
}

func (s *ObligacionesService) BorrarObligacion(id uint) error {
	if !s.obligaciones.BorrarObligacion(id) {
		return ErrNoEncontrado
	}
	return nil
}
