package handlers

import (
	"condominio-api/internal/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// ─── NOTIFICACIONES ──────────────────────────────────────────────────────────

// ListarNotificaciones maneja GET /api/v1/notificaciones
func (s *Server) ListarNotificaciones(w http.ResponseWriter, r *http.Request) {
	notificaciones := s.Notificacion.ListarNotificaciones() // ← Usa NotificacionService
	RespondJSON(w, http.StatusOK, notificaciones)
}

// GetNotificacion maneja GET /api/v1/notificaciones/{id}
func (s *Server) GetNotificacion(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	notif, err := s.Notificacion.ObtenerNotificacion(uint(id))
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, notif)
}

// CreateNotificacion maneja POST /api/v1/notificaciones
func (s *Server) CreateNotificacion(w http.ResponseWriter, r *http.Request) {
	var body models.Notificacion
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}

	creada, err := s.Notificacion.CrearNotificacion(body)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, creada)
}

// MarcarLeida maneja POST /api/v1/notificaciones/{id}/leer
func (s *Server) MarcarLeida(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	notif, err := s.Notificacion.MarcarLeida(uint(id))
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, notif)
}

// DeleteNotificacion maneja DELETE /api/v1/notificaciones/{id}
func (s *Server) DeleteNotificacion(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	if err := s.Notificacion.BorrarNotificacion(uint(id)); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusNoContent, nil)
}
