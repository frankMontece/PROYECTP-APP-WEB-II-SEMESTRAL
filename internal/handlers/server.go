package handlers

import "condominio-api/internal/service"

// Server agrupa todos los servicios que los handlers necesitan
// Siguiendo el patrón de Winter: cada entidad tiene su propio servicio
type Server struct {
	Auth         *service.AuthService
	Area         *service.AreaService
	Reserva      *service.ReservaService
	Notificacion *service.NotificacionService
}

// NewServer crea un nuevo servidor con todos los servicios inyectados
func NewServer(
	auth *service.AuthService,
	area *service.AreaService,
	reserva *service.ReservaService,
	notificacion *service.NotificacionService,
) *Server {
	return &Server{
		Auth:         auth,
		Area:         area,
		Reserva:      reserva,
		Notificacion: notificacion,
	}
}
