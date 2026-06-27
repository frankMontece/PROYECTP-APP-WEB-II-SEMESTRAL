package handlers

import (
	"condominio-api/internal/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// ─── RESERVAS ────────────────────────────────────────────────────────────────

// ListarReservas maneja GET /api/v1/reservas-areas
func (s *Server) ListarReservas(w http.ResponseWriter, r *http.Request) {
	reservas := s.Social.ListarReservas()
	RespondJSON(w, http.StatusOK, reservas)
}

// GetReserva maneja GET /api/v1/reservas-areas/{id}
func (s *Server) GetReserva(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	reserva, err := s.Social.ObtenerReserva(uint(id))
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, reserva)
}

// CreateReserva maneja POST /api/v1/reservas-areas
func (s *Server) CreateReserva(w http.ResponseWriter, r *http.Request) {
	var body models.ReservaArea
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}

	creada, err := s.Social.CrearReserva(body)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, creada)
}

// UpdateReserva maneja PUT /api/v1/reservas-areas/{id}
func (s *Server) UpdateReserva(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var body models.ReservaArea
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}

	actualizada, err := s.Social.ActualizarReserva(uint(id), body)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, actualizada)
}

// DeleteReserva maneja DELETE /api/v1/reservas-areas/{id}
func (s *Server) DeleteReserva(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	if err := s.Social.BorrarReserva(uint(id)); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusNoContent, nil)
}
