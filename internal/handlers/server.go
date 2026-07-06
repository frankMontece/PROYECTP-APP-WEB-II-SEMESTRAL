package handlers

import "condominio-api/internal/service"

// Deps agrupa las dependencias obligatorias del servidor.
// Agregar una entidad = un campo + una linea, sin reordenar argumentos.
type Deps struct {
	Auth         *service.AuthService
	Area         *service.AreaService
	Reserva      *service.ReservaService
	Notificacion *service.NotificacionService
}

// Server agrupa todos los servicios que los handlers necesitan.
type Server struct {
	Auth         *service.AuthService
	Area         *service.AreaService
	Reserva      *service.ReservaService
	Notificacion *service.NotificacionService
}

// NewServer crea un nuevo servidor a partir de Deps.
func NewServer(d Deps) *Server {
	return &Server{
		Auth:         d.Auth,
		Area:         d.Area,
		Reserva:      d.Reserva,
		Notificacion: d.Notificacion,
	}
}
