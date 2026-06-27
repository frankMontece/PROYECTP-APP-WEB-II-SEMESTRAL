package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// parseUintID extrae y valida el parámetro {id} de la URL.
// Devuelve (id, true) si es válido o (0, false) si no lo es.
// Se mantiene para compatibilidad con otros módulos que puedan usarlo.
func parseUintID(r *http.Request) (uint, bool) {
	idStr := chi.URLParam(r, "id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, false
	}
	return uint(id64), true
}

// MontarRutasSocial registra todas las rutas del Módulo B en el router dado.
// Versión actualizada para usar el Server con handlers delgados.
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
