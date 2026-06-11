package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"condominio-api/internal/models"
	"condominio-api/internal/storage"
)

// GetAllAreas atiende GET /api/v1/areas-sociales
func GetAllAreas(w http.ResponseWriter, r *http.Request, store storage.AlmacenSocial) {
	areas := store.ListarAreas()
	RespondJSON(w, http.StatusOK, areas)
}

// GetArea atiende GET /api/v1/areas-sociales/{id}
func GetArea(w http.ResponseWriter, r *http.Request, store storage.AlmacenSocial) {
	id, ok := parseUintID(r)
	if !ok {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero positivo")
		return
	}
	area, encontrada := store.BuscarAreaPorID(id)
	if !encontrada {
		RespondError(w, http.StatusNotFound, "área social no encontrada")
		return
	}
	RespondJSON(w, http.StatusOK, area)
}

// CreateArea atiende POST /api/v1/areas-sociales
func CreateArea(w http.ResponseWriter, r *http.Request, store storage.AlmacenSocial) {
	var body models.AreaSocial
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	if strings.TrimSpace(body.Nombre) == "" {
		RespondError(w, http.StatusBadRequest, "nombre es requerido")
		return
	}
	if body.Capacidad <= 0 {
		RespondError(w, http.StatusBadRequest, "capacidad debe ser mayor a 0")
		return
	}
	body.Activo = true
	creada := store.CrearArea(body)
	RespondJSON(w, http.StatusCreated, creada)
}

// UpdateArea atiende PUT /api/v1/areas-sociales/{id}
func UpdateArea(w http.ResponseWriter, r *http.Request, store storage.AlmacenSocial) {
	id, ok := parseUintID(r)
	if !ok {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero positivo")
		return
	}
	var body models.AreaSocial
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	if strings.TrimSpace(body.Nombre) == "" {
		RespondError(w, http.StatusBadRequest, "nombre es requerido")
		return
	}
	actualizada, encontrada := store.ActualizarArea(id, body)
	if !encontrada {
		RespondError(w, http.StatusNotFound, "área social no encontrada")
		return
	}
	RespondJSON(w, http.StatusOK, actualizada)
}

// DeleteArea atiende DELETE /api/v1/areas-sociales/{id}
func DeleteArea(w http.ResponseWriter, r *http.Request, store storage.AlmacenSocial) {
	id, ok := parseUintID(r)
	if !ok {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero positivo")
		return
	}
	if !store.BorrarArea(id) {
		RespondError(w, http.StatusNotFound, "área social no encontrada")
		return
	}
	RespondJSON(w, http.StatusNoContent, nil)
}
