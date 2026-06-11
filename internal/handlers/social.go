package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"condominio-api/internal/storage"
)

// parseUintID extrae y valida el parámetro {id} de la URL.
// Devuelve (id, true) si es válido o (0, false) si no lo es.
func parseUintID(r *http.Request) (uint, bool) {
	idStr := chi.URLParam(r, "id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, false
	}
	return uint(id64), true
}

// MontarRutasSocial registra todas las rutas del Módulo B en el router dado.
// Se llama desde main.go dentro del bloque r.Route("/api/v1", ...).
func MontarRutasSocial(r chi.Router, store storage.AlmacenSocial) {

	// ── Áreas Sociales ─────────────────────────────────────────────────────
	r.Get("/areas-sociales", func(w http.ResponseWriter, req *http.Request) {
		GetAllAreas(w, req, store)
	})
	r.Post("/areas-sociales", func(w http.ResponseWriter, req *http.Request) {
		CreateArea(w, req, store)
	})
	r.Get("/areas-sociales/{id}", func(w http.ResponseWriter, req *http.Request) {
		GetArea(w, req, store)
	})
	r.Put("/areas-sociales/{id}", func(w http.ResponseWriter, req *http.Request) {
		UpdateArea(w, req, store)
	})
	r.Delete("/areas-sociales/{id}", func(w http.ResponseWriter, req *http.Request) {
		DeleteArea(w, req, store)
	})

	// ── Reservas de Áreas ──────────────────────────────────────────────────
	r.Get("/reservas-areas", func(w http.ResponseWriter, req *http.Request) {
		GetAllReservas(w, req, store)
	})
	r.Post("/reservas-areas", func(w http.ResponseWriter, req *http.Request) {
		CreateReserva(w, req, store)
	})
	r.Get("/reservas-areas/{id}", func(w http.ResponseWriter, req *http.Request) {
		GetReserva(w, req, store)
	})
	r.Put("/reservas-areas/{id}", func(w http.ResponseWriter, req *http.Request) {
		UpdateReserva(w, req, store)
	})
	r.Delete("/reservas-areas/{id}", func(w http.ResponseWriter, req *http.Request) {
		DeleteReserva(w, req, store)
	})

	// ── Notificaciones ─────────────────────────────────────────────────────
	r.Get("/notificaciones", func(w http.ResponseWriter, req *http.Request) {
		GetAllNotificaciones(w, req, store)
	})
	r.Post("/notificaciones", func(w http.ResponseWriter, req *http.Request) {
		CreateNotificacion(w, req, store)
	})
	r.Get("/notificaciones/{id}", func(w http.ResponseWriter, req *http.Request) {
		GetNotificacion(w, req, store)
	})
	r.Post("/notificaciones/{id}/leer", func(w http.ResponseWriter, req *http.Request) {
		MarcarNotificacionLeida(w, req, store)
	})
	r.Delete("/notificaciones/{id}", func(w http.ResponseWriter, req *http.Request) {
		DeleteNotificacion(w, req, store)
	})
}
