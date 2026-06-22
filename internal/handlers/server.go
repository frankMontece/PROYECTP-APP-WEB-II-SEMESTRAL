package handlers

import "condominio-api/internal/service"

// Server agrupa los servicios que los handlers necesitan
type Server struct {
	Social *service.SocialService
	Auth   *service.AuthService
}

// NewServer crea un nuevo servidor con los servicios inyectados
func NewServer(social *service.SocialService, auth *service.AuthService) *Server {
	return &Server{Social: social, Auth: auth}
}
