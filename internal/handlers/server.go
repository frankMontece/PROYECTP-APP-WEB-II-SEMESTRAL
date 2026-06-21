package handlers

import "condominio-api/internal/service"

// Server agrupa todos los servicios del módulo.
// Los handlers son métodos de este struct, igual que en cafetería.
type Server struct {
	Auth      *service.AuthService
	Vehiculos *service.VehiculoService
	Visitas   *service.VisitaService
	Accesos   *service.AccesoService
}

func NewServer(
	auth *service.AuthService,
	vehiculos *service.VehiculoService,
	visitas *service.VisitaService,
	accesos *service.AccesoService,
) *Server {
	return &Server{
		Auth:      auth,
		Vehiculos: vehiculos,
		Visitas:   visitas,
		Accesos:   accesos,
	}
}
