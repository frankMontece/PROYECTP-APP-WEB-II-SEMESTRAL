package handlers

import (
	"condominio-api/internal/service"
	"encoding/json"
	"errors"
	"log"
	"net/http"
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

// statusDeError mapea errores de dominio a códigos HTTP
func statusDeError(err error) int {
	switch {
	case errors.Is(err, service.ErrNoEncontrado):
		return http.StatusNotFound
	case errors.Is(err, service.ErrNombreVacio),
		errors.Is(err, service.ErrCapacidadInvalida),
		errors.Is(err, service.ErrFechasInvalidas),
		errors.Is(err, service.ErrNumeroPersonasInvalido):
		return http.StatusBadRequest
	case errors.Is(err, service.ErrEmailEnUso):
		return http.StatusConflict
	case errors.Is(err, service.ErrCredencialesInvalidas):
		return http.StatusUnauthorized
	case errors.Is(err, service.ErrSinPermiso):
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
