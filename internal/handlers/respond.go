package handlers

import (
	"condominio-api/internal/service"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

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

func RespondError(w http.ResponseWriter, status int, mensaje string) {
	RespondJSON(w, status, map[string]string{"error": mensaje})
}

// statusDeError convierte un error de servicio en el HTTP status code correcto.
// Centraliza el mapeo para que ningún handler tenga que hardcodear status codes.
func statusDeError(err error) int {
	switch {
	case errors.Is(err, service.ErrCredencialesInvalidas):
		return http.StatusUnauthorized
	case errors.Is(err, service.ErrEmailEnUso):
		return http.StatusConflict
	case errors.Is(err, service.ErrVehiculoNoEncontrado),
		errors.Is(err, service.ErrVisitaNoEncontrada),
		errors.Is(err, service.ErrAccesoNoEncontrado):
		return http.StatusNotFound
	case errors.Is(err, service.ErrPlacaRequerida),
		errors.Is(err, service.ErrMarcaRequerida),
		errors.Is(err, service.ErrResidenteRequerido),
		errors.Is(err, service.ErrPlacaVisitanteRequerida),
		errors.Is(err, service.ErrNombreVisitanteRequerido),
		errors.Is(err, service.ErrCondominioRequerido),
		errors.Is(err, service.ErrVehiculoRequerido),
		errors.Is(err, service.ErrTipoMovimientoInvalido):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
