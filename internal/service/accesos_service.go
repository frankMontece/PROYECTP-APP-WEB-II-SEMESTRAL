package service

import (
	"condominio-api/internal/models"
	"condominio-api/internal/storage"
	"time"
)

// VehiculoService contiene la lógica de negocio para Vehículos
type AccesosService struct {
	repo storage.AccesosRepository
}

// NewVehiculoService crea una nueva instancia del servicio de vehículos
func NewAccesosService(repo storage.AccesosRepository) *AccesosService {
	return &AccesosService{repo: repo}
}

// ACCESO SERVICE

type AccesoService struct {
	repo storage.AccesosRepository
}

func NewAccesoService(repo storage.AccesosRepository) *AccesoService {
	return &AccesoService{repo: repo}
}

func (s *AccesoService) Listar() []models.AccesoVehiculo {
	return s.repo.ListarAccesos()
}

func (s *AccesoService) Obtener(id uint) (models.AccesoVehiculo, error) {
	a, ok := s.repo.BuscarAccesoPorID(id)
	if !ok {
		return models.AccesoVehiculo{}, ErrAccesoNoEncontrado
	}
	return a, nil
}

func (s *AccesoService) Crear(a models.AccesoVehiculo) (models.AccesoVehiculo, error) {
	if err := validarAcceso(a); err != nil {
		return models.AccesoVehiculo{}, err
	}
	// Lógica de negocio: el servidor registra la fecha/hora, no el cliente
	a.FechaHora = time.Now()
	return s.repo.CrearAcceso(a), nil
}

func (s *AccesoService) Borrar(id uint) error {
	if !s.repo.BorrarAcceso(id) {
		return ErrAccesoNoEncontrado
	}
	return nil
}

func validarAcceso(a models.AccesoVehiculo) error {
	if a.VehiculoID == 0 {
		return ErrVehiculoRequerido
	}
	if a.CondominioID == 0 {
		return ErrCondominioRequerido
	}
	if a.TipoMovimiento != "entrada" && a.TipoMovimiento != "salida" {
		return ErrTipoMovimientoInvalido
	}
	return nil
}
