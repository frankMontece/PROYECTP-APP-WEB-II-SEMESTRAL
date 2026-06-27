package handlers

import (
	"encoding/json"
	"net/http"

	"condominio-api/internal/models"
)

func (s *Server) ListarAccesos(w http.ResponseWriter, r *http.Request) {
	RespondJSON(w, http.StatusOK, s.Accesos.Listar())
}

func (s *Server) ObtenerAcceso(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero positivo")
		return
	}

	acceso, err := s.Accesos.Obtener(id)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, acceso)
}

func (s *Server) CrearAcceso(w http.ResponseWriter, r *http.Request) {
	var body models.AccesoVehiculo
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}

	acceso, err := s.Accesos.Crear(body)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusCreated, acceso)
}

func (s *Server) BorrarAcceso(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero positivo")
		return
	}

	if err := s.Accesos.Borrar(id); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusNoContent, nil)
}
