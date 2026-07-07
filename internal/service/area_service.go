package service

import (
	"condominio-api/internal/models"
	"condominio-api/internal/storage"
	"strings"
)

// AreaService contiene la lógica de negocio para Áreas Sociales
type AreaService struct {
	repo storage.AreaSocialRepository
}

// NewAreaService crea una nueva instancia del servicio de áreas
func NewAreaService(repo storage.AreaSocialRepository) *AreaService {
	return &AreaService{repo: repo}
}

// ListarAreas retorna todas las áreas
func (s *AreaService) ListarAreas() []models.AreaSocial {
	return s.repo.ListarAreas()
}

// ObtenerArea retorna un área por ID
func (s *AreaService) ObtenerArea(id uint) (models.AreaSocial, error) {
	a, ok := s.repo.BuscarAreaPorID(id)
	if !ok {
		return models.AreaSocial{}, ErrNoEncontrado
	}
	return a, nil
}

// CrearArea valida y crea un área social
func (s *AreaService) CrearArea(a models.AreaSocial) (models.AreaSocial, error) {
	if strings.TrimSpace(a.Nombre) == "" {
		return models.AreaSocial{}, ErrNombreVacio
	}
	if a.Capacidad <= 0 {
		return models.AreaSocial{}, ErrCapacidadInvalida
	}
	a.Activo = true
	return s.repo.CrearArea(a), nil
}

// ActualizarArea valida y actualiza un área
func (s *AreaService) ActualizarArea(id uint, datos models.AreaSocial) (models.AreaSocial, error) {
	if strings.TrimSpace(datos.Nombre) == "" {
		return models.AreaSocial{}, ErrNombreVacio
	}
	a, ok := s.repo.ActualizarArea(id, datos)
	if !ok {
		return models.AreaSocial{}, ErrNoEncontrado
	}
	return a, nil
}

// BorrarArea elimina un área
func (s *AreaService) BorrarArea(id uint) error {
	if !s.repo.BorrarArea(id) {
		return ErrNoEncontrado
	}
	return nil
}
