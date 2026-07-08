package service

import (
	"condominio-api/internal/models"
	"condominio-api/internal/storage"
	"strings"
	"time"
)

// VehiculoService contiene la lógica de negocio para Vehículos
type VehiculoService struct {
	repo storage.VehiculosRepository
}

// NewVehiculoService crea una nueva instancia del servicio de vehículos
func NewVehiculoService(repo storage.VehiculosRepository) *VehiculoService {
	return &VehiculoService{repo: repo}
}

func (s *VehiculoService) Listar() []models.Vehiculo {
	return s.repo.ListarVehiculos()
}

func (s *VehiculoService) Obtener(id uint) (models.Vehiculo, error) {
	v, ok := s.repo.BuscarVehiculoPorID(id)
	if !ok {
		return models.Vehiculo{}, ErrVehiculoNoEncontrado
	}
	return v, nil
}

func (s *VehiculoService) Crear(v models.Vehiculo) (models.Vehiculo, error) {
	if err := validarVehiculo(v); err != nil {
		return models.Vehiculo{}, err
	}
	// Lógica de negocio: el servidor asigna estos campos, no el cliente
	v.Activo = true
	v.CreatedAt = time.Now()
	return s.repo.CrearVehiculo(v), nil
}

func (s *VehiculoService) Actualizar(id uint, v models.Vehiculo) (models.Vehiculo, error) {
	if err := validarVehiculo(v); err != nil {
		return models.Vehiculo{}, err
	}
	v, ok := s.repo.ActualizarVehiculo(id, v)
	if !ok {
		return models.Vehiculo{}, ErrVehiculoNoEncontrado
	}
	return v, nil
}

func (s *VehiculoService) Borrar(id uint) error {
	if !s.repo.BorrarVehiculo(id) {
		return ErrVehiculoNoEncontrado
	}
	return nil
}

func validarVehiculo(v models.Vehiculo) error {
	if v.ResidenteID == 0 {
		return ErrResidenteRequerido
	}
	if strings.TrimSpace(v.Placa) == "" {
		return ErrPlacaRequerida
	}
	if strings.TrimSpace(v.Marca) == "" {
		return ErrMarcaRequerida
	}
	return nil
}
