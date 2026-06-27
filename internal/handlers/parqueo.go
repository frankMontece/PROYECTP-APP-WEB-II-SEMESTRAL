package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"condominio-api/internal/models"
)

// parseID extrae el parámetro "id" de la URL y lo convierte a uint.
func parseID(r *http.Request) (uint, error) {
	id64, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id64), nil
}

// =================================================================
// VEHICULOS
// =================================================================

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

// =================================================================
// VISITAS
// =================================================================

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

// =================================================================
// ACCESOS
// =================================================================

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
