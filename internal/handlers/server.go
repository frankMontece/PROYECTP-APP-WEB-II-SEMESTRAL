package handlers

import "condominio-api/internal/service"

// Services agrupa todos los servicios de los tres módulos.
// Agregar un dominio nuevo = un campo aquí, sin tocar el constructor.
type Services struct {
	// Compartido
	Auth *service.AuthService

	// Módulo Obligaciones
	Obligaciones *service.ObligacionesService
	Multas       *service.MultasService

	// Módulo Parqueo (acceso y vehículos)
	//	Vehiculos *service.VehiculoService
	//	Visitas   *service.VisitaService
	//	Accesos   *service.AccesoService

	// Módulo Área Social
	Area         *service.AreaService
	Reserva      *service.ReservaService
	Notificacion *service.NotificacionService
}

// Server agrupa todos los servicios que los handlers necesitan,
// vía embedding — permite escribir server.Auth, server.Vehiculos, etc.
type Server struct {
	Services
}

// NewServer crea un nuevo servidor a partir de Services.
func NewServer(s Services) *Server {
	return &Server{
		Services: s,
	}
}
