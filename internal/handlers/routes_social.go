package handlers

import (
	"github.com/go-chi/chi/v5"

	"condominio-api/internal/middleware"
)

// MontarRutasSocial registra todas las rutas del Módulo B (Área Social,
// Reservas, Notificaciones). Lectura (GET) abierta a cualquier usuario
// autenticado; escritura (POST/PUT/DELETE) restringida a rol "admin".
func MontarRutasSocial(r chi.Router, s *Server) {
	soloAdmin := middleware.RequireRol("admin")

	// ── Áreas Sociales ─────────────────────────────────────────────────────
	r.Get("/areas-sociales", s.ListarAreas)
	r.Get("/areas-sociales/{id}", s.GetArea)
	r.With(soloAdmin).Post("/areas-sociales", s.CreateArea)
	r.With(soloAdmin).Put("/areas-sociales/{id}", s.UpdateArea)
	r.With(soloAdmin).Delete("/areas-sociales/{id}", s.DeleteArea)

	// ── Reservas de Áreas ──────────────────────────────────────────────────
	// Nota: crear una reserva es una acción típica de residente, no solo
	// de admin — se deja abierta a cualquier autenticado. Solo actualizar
	// (aprobar/rechazar) y borrar quedan restringidos a admin.
	r.Get("/reservas-areas", s.ListarReservas)
	r.Get("/reservas-areas/{id}", s.GetReserva)
	r.Post("/reservas-areas", s.CreateReserva)
	r.With(soloAdmin).Put("/reservas-areas/{id}", s.UpdateReserva)
	r.With(soloAdmin).Delete("/reservas-areas/{id}", s.DeleteReserva)

	// ── Notificaciones ─────────────────────────────────────────────────────
	r.Get("/notificaciones", s.ListarNotificaciones)
	r.Get("/notificaciones/{id}", s.GetNotificacion)
	r.Post("/notificaciones/{id}/leer", s.MarcarLeida)
	r.With(soloAdmin).Post("/notificaciones", s.CreateNotificacion)
	r.With(soloAdmin).Delete("/notificaciones/{id}", s.DeleteNotificacion)
}
