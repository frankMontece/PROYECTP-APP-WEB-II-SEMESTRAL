package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"condominio-api/internal/models"
	"condominio-api/internal/storage"
)

// GetAllReservas atiende GET /api/v1/reservas-areas
func GetAllReservas(w http.ResponseWriter, r *http.Request, store storage.AlmacenSocial) {
	reservas := store.ListarReservas()
	RespondJSON(w, http.StatusOK, reservas)
}

// GetReserva atiende GET /api/v1/reservas-areas/{id}
func GetReserva(w http.ResponseWriter, r *http.Request, store storage.AlmacenSocial) {
	id, ok := parseUintID(r)
	if !ok {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero positivo")
		return
	}
	reserva, encontrada := store.BuscarReservaPorID(id)
	if !encontrada {
		RespondError(w, http.StatusNotFound, "reserva no encontrada")
		return
	}
	RespondJSON(w, http.StatusOK, reserva)
}

// CreateReserva atiende POST /api/v1/reservas-areas
func CreateReserva(w http.ResponseWriter, r *http.Request, store storage.AlmacenSocial) {
	var body models.ReservaArea
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	if body.AreaID == 0 {
		RespondError(w, http.StatusBadRequest, "area_id es requerido")
		return
	}
	if body.ResidenteID == 0 {
		RespondError(w, http.StatusBadRequest, "residente_id es requerido")
		return
	}
	if strings.TrimSpace(body.Proposito) == "" {
		RespondError(w, http.StatusBadRequest, "proposito es requerido")
		return
	}
	if body.NumeroPersonas <= 0 {
		RespondError(w, http.StatusBadRequest, "numero_personas debe ser mayor a 0")
		return
	}
	if body.FechaInicio.IsZero() || body.FechaFin.IsZero() {
		RespondError(w, http.StatusBadRequest, "fecha_inicio y fecha_fin son requeridas")
		return
	}
	if !body.FechaFin.After(body.FechaInicio) {
		RespondError(w, http.StatusBadRequest, "fecha_fin debe ser posterior a fecha_inicio")
		return
	}
	creada := store.CrearReserva(body)
	RespondJSON(w, http.StatusCreated, creada)
}

// UpdateReserva atiende PUT /api/v1/reservas-areas/{id}
func UpdateReserva(w http.ResponseWriter, r *http.Request, store storage.AlmacenSocial) {
	id, ok := parseUintID(r)
	if !ok {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero positivo")
		return
	}
	var body models.ReservaArea
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	if !body.FechaInicio.IsZero() && !body.FechaFin.IsZero() {
		if !body.FechaFin.After(body.FechaInicio) {
			RespondError(w, http.StatusBadRequest, "fecha_fin debe ser posterior a fecha_inicio")
			return
		}
	}
	actualizada, encontrada := store.ActualizarReserva(id, body)
	if !encontrada {
		RespondError(w, http.StatusNotFound, "reserva no encontrada")
		return
	}
	RespondJSON(w, http.StatusOK, actualizada)
}

// DeleteReserva atiende DELETE /api/v1/reservas-areas/{id}
func DeleteReserva(w http.ResponseWriter, r *http.Request, store storage.AlmacenSocial) {
	id, ok := parseUintID(r)
	if !ok {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero positivo")
		return
	}
	if !store.BorrarReserva(id) {
		RespondError(w, http.StatusNotFound, "reserva no encontrada")
		return
	}
	RespondJSON(w, http.StatusNoContent, nil)
}
