package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"condominio-api/internal/models"
	"condominio-api/internal/storage"
)

// GetAllNotificaciones atiende GET /api/v1/notificaciones
func GetAllNotificaciones(w http.ResponseWriter, r *http.Request, store storage.AlmacenSocial) {
	notificaciones := store.ListarNotificaciones()
	RespondJSON(w, http.StatusOK, notificaciones)
}

// GetNotificacion atiende GET /api/v1/notificaciones/{id}
func GetNotificacion(w http.ResponseWriter, r *http.Request, store storage.AlmacenSocial) {
	id, ok := parseUintID(r)
	if !ok {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero positivo")
		return
	}
	notif, encontrada := store.BuscarNotificacionPorID(id)
	if !encontrada {
		RespondError(w, http.StatusNotFound, "notificación no encontrada")
		return
	}
	RespondJSON(w, http.StatusOK, notif)
}

// CreateNotificacion atiende POST /api/v1/notificaciones
func CreateNotificacion(w http.ResponseWriter, r *http.Request, store storage.AlmacenSocial) {
	var body models.Notificacion
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	if body.ResidenteID == 0 {
		RespondError(w, http.StatusBadRequest, "residente_id es requerido")
		return
	}
	if strings.TrimSpace(body.Tipo) == "" {
		RespondError(w, http.StatusBadRequest, "tipo es requerido")
		return
	}
	if strings.TrimSpace(body.Mensaje) == "" {
		RespondError(w, http.StatusBadRequest, "mensaje es requerido")
		return
	}
	creada := store.CrearNotificacion(body)
	RespondJSON(w, http.StatusCreated, creada)
}

// MarcarNotificacionLeida atiende POST /api/v1/notificaciones/{id}/leer
func MarcarNotificacionLeida(w http.ResponseWriter, r *http.Request, store storage.AlmacenSocial) {
	id, ok := parseUintID(r)
	if !ok {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero positivo")
		return
	}
	notif, encontrada := store.MarcarComoLeida(id)
	if !encontrada {
		RespondError(w, http.StatusNotFound, "notificación no encontrada")
		return
	}
	RespondJSON(w, http.StatusOK, notif)
}

// DeleteNotificacion atiende DELETE /api/v1/notificaciones/{id}
func DeleteNotificacion(w http.ResponseWriter, r *http.Request, store storage.AlmacenSocial) {
	id, ok := parseUintID(r)
	if !ok {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero positivo")
		return
	}
	if !store.BorrarNotificacion(id) {
		RespondError(w, http.StatusNotFound, "notificación no encontrada")
		return
	}
	RespondJSON(w, http.StatusNoContent, nil)
}
