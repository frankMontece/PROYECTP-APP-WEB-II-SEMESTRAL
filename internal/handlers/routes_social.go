package handlers

import (
	"github.com/go-chi/chi/v5"
)

// MontarRutasSocial registra todas las rutas del Módulo B
// Usa los handlers que llaman a servicios específicos
func MontarRutasSocial(r chi.Router, s *Server) {
	// ── Áreas Sociales ─────────────────────────────────────────────────────
	r.Get("/areas-sociales", s.ListarAreas)
	r.Post("/areas-sociales", s.CreateArea)
	r.Get("/areas-sociales/{id}", s.GetArea)
	r.Put("/areas-sociales/{id}", s.UpdateArea)
	r.Delete("/areas-sociales/{id}", s.DeleteArea)

	// ── Reservas de Áreas ──────────────────────────────────────────────────
	r.Get("/reservas-areas", s.ListarReservas)
	r.Post("/reservas-areas", s.CreateReserva)
	r.Get("/reservas-areas/{id}", s.GetReserva)
	r.Put("/reservas-areas/{id}", s.UpdateReserva)
	r.Delete("/reservas-areas/{id}", s.DeleteReserva)

	// ── Notificaciones ─────────────────────────────────────────────────────
	r.Get("/notificaciones", s.ListarNotificaciones)
	r.Post("/notificaciones", s.CreateNotificacion)
	r.Get("/notificaciones/{id}", s.GetNotificacion)
	r.Post("/notificaciones/{id}/leer", s.MarcarLeida)
	r.Delete("/notificaciones/{id}", s.DeleteNotificacion)
}
