package handlers

import (
	"encoding/json"
	"net/http"

	"condominio-api/internal/models"
)

// ─── Obligaciones ─────────────────────────────────────────────────────────────

func (s *Server) ListarObligaciones(w http.ResponseWriter, r *http.Request) {
	RespondJSON(w, http.StatusOK, s.Obligaciones.ListarObligaciones())
}

func (s *Server) GetObligacion(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID debe ser un número entero positivo")
		return
	}
	o, err := s.Obligaciones.ObtenerObligacion(id)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, o)
}

func (s *Server) CreateObligacion(w http.ResponseWriter, r *http.Request) {
	var body models.Obligacion
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	creada, err := s.Obligaciones.CrearObligacion(body)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, creada)
}

func (s *Server) UpdateObligacion(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID debe ser un número entero positivo")
		return
	}
	var body models.Obligacion
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	actualizada, err := s.Obligaciones.ActualizarObligacion(id, body)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, actualizada)
}

func (s *Server) DeleteObligacion(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID debe ser un número entero positivo")
		return
	}
	if err := s.Obligaciones.BorrarObligacion(id); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, map[string]string{"mensaje": "obligacion eliminada"})
}
