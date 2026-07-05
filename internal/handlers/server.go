package handlers

import "condominio-api/internal/service"

// Server agrupa todos los servicios del módulo.
// Los handlers son métodos de este struct, igual que en cafetería.
type Services struct {
	Auth      *service.AuthService
	Vehiculos *service.VehiculoService
	Visitas   *service.VisitaService
	Accesos   *service.AccesoService
}

type Server struct {
	Services
}

func NewServer(s Services) *Server {
	return &Server{
		Services: s,
	}
}
