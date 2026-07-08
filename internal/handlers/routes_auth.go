package handlers

import "github.com/go-chi/chi/v5"

func (s *Server) MontarRutasAuth(r chi.Router) {
	r.Post("/auth/register", s.Registrar)
	r.Post("/auth/login", s.Login)
}
