package service

import (
	"fmt"
	"strings"
	"time"

	"condominio-api/internal/models"
	"condominio-api/internal/storage"
)

// VEHICULO SERVICE

type VehiculoService struct {
	repo storage.VehiculosRepository
}

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
