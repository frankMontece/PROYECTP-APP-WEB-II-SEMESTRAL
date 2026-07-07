package handlers

import "condominio-api/internal/service"

// Server agrupa todos los servicios del módulo A.
type Server struct {
	Auth         *service.AuthService
	Obligaciones *service.ObligacionesService
	Multas       *service.MultasService
}

func NewServer(
	auth *service.AuthService,
	obligaciones *service.ObligacionesService,
	multas *service.MultasService,
) *Server {
	return &Server{
		Auth:         auth,
		Obligaciones: obligaciones,
		Multas:       multas,
	}
}
