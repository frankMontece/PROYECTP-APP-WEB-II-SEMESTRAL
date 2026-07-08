package handlers

import (
	"encoding/json"
	"net/http"

	"condominio-api/internal/models"
)

func (s *Server) ListarMultas(w http.ResponseWriter, r *http.Request) {
	RespondJSON(w, http.StatusOK, s.Multas.ListarMultas())
}

func (s *Server) GetMulta(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID debe ser un número entero positivo")
		return
	}
	m, err := s.Multas.ObtenerMulta(id)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, m)
}

func (s *Server) CreateMulta(w http.ResponseWriter, r *http.Request) {
	var body models.Multa
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	creada, err := s.Multas.CrearMulta(body)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, creada)
}

func (s *Server) UpdateMulta(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID debe ser un número entero positivo")
		return
	}
	var body models.Multa
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	actualizada, err := s.Multas.ActualizarMulta(id, body)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, actualizada)
}

func (s *Server) DeleteMulta(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID debe ser un número entero positivo")
		return
	}
	if err := s.Multas.BorrarMulta(id); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, map[string]string{"mensaje": "multa eliminada"})
}
