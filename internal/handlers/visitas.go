package handlers

import (
	"encoding/json"
	"net/http"

	"condominio-api/internal/models"
)

func (s *Server) ListarVisitas(w http.ResponseWriter, r *http.Request) {
	RespondJSON(w, http.StatusOK, s.Visitas.Listar())
}

func (s *Server) ObtenerVisita(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero positivo")
		return
	}

	visita, err := s.Visitas.Obtener(id)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, visita)
}

func (s *Server) CrearVisita(w http.ResponseWriter, r *http.Request) {
	var body models.VisitaVehiculo
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}

	visita, err := s.Visitas.Crear(body)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusCreated, visita)
}

func (s *Server) ActualizarVisita(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero positivo")
		return
	}

	var body models.VisitaVehiculo
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}

	visita, err := s.Visitas.Actualizar(id, body)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, visita)
}

// RegistrarSalida atiende POST /api/v1/visitas/{id}/salida.
func (s *Server) RegistrarSalida(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero positivo")
		return
	}

	visita, err := s.Visitas.RegistrarSalida(id)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, visita)
}

func (s *Server) BorrarVisita(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero positivo")
		return
	}

	if err := s.Visitas.Borrar(id); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusNoContent, nil)
}
