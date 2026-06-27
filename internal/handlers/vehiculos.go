package handlers

import (
	"encoding/json"
	"net/http"

	"condominio-api/internal/models"
)

func (s *Server) ListarVehiculos(w http.ResponseWriter, r *http.Request) {
	RespondJSON(w, http.StatusOK, s.Vehiculos.Listar())
}

func (s *Server) ObtenerVehiculo(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero positivo")
		return
	}

	vehiculo, err := s.Vehiculos.Obtener(id)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, vehiculo)
}

func (s *Server) CrearVehiculo(w http.ResponseWriter, r *http.Request) {
	var body models.Vehiculo
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}

	vehiculo, err := s.Vehiculos.Crear(body)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusCreated, vehiculo)
}

func (s *Server) ActualizarVehiculo(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero positivo")
		return
	}

	var body models.Vehiculo
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}

	vehiculo, err := s.Vehiculos.Actualizar(id, body)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, vehiculo)
}

func (s *Server) BorrarVehiculo(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero positivo")
		return
	}

	if err := s.Vehiculos.Borrar(id); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusNoContent, nil)
}
