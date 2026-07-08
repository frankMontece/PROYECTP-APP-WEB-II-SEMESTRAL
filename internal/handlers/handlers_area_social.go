package handlers

import (
	"condominio-api/internal/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// ─── ÁREAS SOCIALES ──────────────────────────────────────────────────────────

// ListarAreas maneja GET /api/v1/areas-sociales
func (s *Server) ListarAreas(w http.ResponseWriter, r *http.Request) {
	areas := s.Area.ListarAreas() // ← Usa AreaService
	RespondJSON(w, http.StatusOK, areas)
}

// GetArea maneja GET /api/v1/areas-sociales/{id}
func (s *Server) GetArea(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	area, err := s.Area.ObtenerArea(uint(id))
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, area)
}

// CreateArea maneja POST /api/v1/areas-sociales
func (s *Server) CreateArea(w http.ResponseWriter, r *http.Request) {
	var body models.AreaSocial
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}

	creada, err := s.Area.CrearArea(body)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, creada)
}

// UpdateArea maneja PUT /api/v1/areas-sociales/{id}
func (s *Server) UpdateArea(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var body models.AreaSocial
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}

	actualizada, err := s.Area.ActualizarArea(uint(id), body)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, actualizada)
}

// DeleteArea maneja DELETE /api/v1/areas-sociales/{id}
func (s *Server) DeleteArea(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	if err := s.Area.BorrarArea(uint(id)); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusNoContent, nil)
}
