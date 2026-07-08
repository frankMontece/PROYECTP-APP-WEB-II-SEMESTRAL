package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"condominio-api/internal/service"
)

// RespondJSON envía una respuesta JSON al cliente
func RespondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("error codificando JSON: %v", err)
	}
}

// RespondError envía una respuesta de error en formato JSON
func RespondError(w http.ResponseWriter, status int, mensaje string) {
	RespondJSON(w, status, map[string]string{"error": mensaje})
}

// statusDeError convierte un error de servicio en el HTTP status code correcto.
// Centraliza el mapeo para que ningún handler tenga que hardcodear status codes.
func statusDeError(err error) int {
	switch {
	// ── 404 Not Found ──────────────────────────────────────────────────
	case errors.Is(err, service.ErrNoEncontrado),
		errors.Is(err, service.ErrVehiculoNoEncontrado),
		errors.Is(err, service.ErrVisitaNoEncontrada),
		errors.Is(err, service.ErrAccesoNoEncontrado):
		return http.StatusNotFound

	// ── 400 Bad Request (validación de campos) ────────────────────────
	case errors.Is(err, service.ErrNombreVacio),
		errors.Is(err, service.ErrCapacidadInvalida),
		errors.Is(err, service.ErrFechasInvalidas),
		errors.Is(err, service.ErrNumeroPersonasInvalido),
		errors.Is(err, service.ErrResidenteIDInvalido),
		errors.Is(err, service.ErrMontoInvalido),
		errors.Is(err, service.ErrPeriodoVacio),
		errors.Is(err, service.ErrMotivoVacio),
		errors.Is(err, service.ErrTipoObligacion),
		errors.Is(err, service.ErrPlacaRequerida),
		errors.Is(err, service.ErrMarcaRequerida),
		errors.Is(err, service.ErrResidenteRequerido),
		errors.Is(err, service.ErrPlacaVisitanteRequerida),
		errors.Is(err, service.ErrNombreVisitanteRequerido),
		errors.Is(err, service.ErrCondominioRequerido),
		errors.Is(err, service.ErrVehiculoRequerido),
		errors.Is(err, service.ErrTipoMovimientoInvalido):
		return http.StatusBadRequest

	// ── 409 Conflict ────────────────────────────────────────────────────
	case errors.Is(err, service.ErrEmailEnUso):
		return http.StatusConflict

	// ── 401 Unauthorized ───────────────────────────────────────────────
	case errors.Is(err, service.ErrCredencialesInvalidas):
		return http.StatusUnauthorized

	// ── 403 Forbidden ──────────────────────────────────────────────────
	case errors.Is(err, service.ErrSinPermiso):
		return http.StatusForbidden

	default:
		return http.StatusInternalServerError
	}
}
