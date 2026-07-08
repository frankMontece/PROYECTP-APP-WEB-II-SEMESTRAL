package service

import (
	"condominio-api/internal/models"
	"condominio-api/internal/storage"
	"fmt"
	"strings"
	"time"
)

// VehiculoService contiene la lógica de negocio para Vehículos
type VisitasService struct {
	repo storage.VisitasRepository
}

// NewVehiculoService crea una nueva instancia del servicio de vehículos
func NewVisitasService(repo storage.VisitasRepository) *VisitaService {
	return &VisitaService{repo: repo}
}

// VISITA SERVICE

type VisitaService struct {
	repo storage.VisitasRepository
}

func NewVisitaService(repo storage.VisitasRepository) *VisitaService {
	return &VisitaService{repo: repo}
}

func (s *VisitaService) Listar() []models.VisitaVehiculo {
	return s.repo.ListarVisitas()
}

func (s *VisitaService) Obtener(id uint) (models.VisitaVehiculo, error) {
	vis, ok := s.repo.BuscarVisitaPorID(id)
	if !ok {
		return models.VisitaVehiculo{}, ErrVisitaNoEncontrada
	}
	return vis, nil
}

func (s *VisitaService) Crear(vis models.VisitaVehiculo) (models.VisitaVehiculo, error) {
	if err := validarVisita(vis); err != nil {
		return models.VisitaVehiculo{}, err
	}
	// Lógica de negocio: el servidor genera el QR y registra la entrada
	vis.CodigoQR = fmt.Sprintf("QR-%d", time.Now().UnixNano())
	vis.EstadoQR = "pendiente"
	ahora := time.Now()
	vis.HoraEntrada = &ahora
	vis.HoraSalida = nil
	return s.repo.CrearVisita(vis), nil
}

func (s *VisitaService) Actualizar(id uint, vis models.VisitaVehiculo) (models.VisitaVehiculo, error) {
	if err := validarVisita(vis); err != nil {
		return models.VisitaVehiculo{}, err
	}
	vis, ok := s.repo.ActualizarVisita(id, vis)
	if !ok {
		return models.VisitaVehiculo{}, ErrVisitaNoEncontrada
	}
	return vis, nil
}

func (s *VisitaService) RegistrarSalida(id uint) (models.VisitaVehiculo, error) {
	vis, ok := s.repo.RegistrarSalidaVisita(id)
	if !ok {
		return models.VisitaVehiculo{}, ErrVisitaNoEncontrada
	}
	return vis, nil
}

func (s *VisitaService) Borrar(id uint) error {
	if !s.repo.BorrarVisita(id) {
		return ErrVisitaNoEncontrada
	}
	return nil
}

func validarVisita(vis models.VisitaVehiculo) error {
	if vis.CondominioID == 0 {
		return ErrCondominioRequerido
	}
	if vis.ResidenteID == 0 {
		return ErrResidenteRequerido
	}
	if strings.TrimSpace(vis.PlacaVisitante) == "" {
		return ErrPlacaVisitanteRequerida
	}
	if strings.TrimSpace(vis.NombreVisitante) == "" {
		return ErrNombreVisitanteRequerido
	}
	return nil
}
