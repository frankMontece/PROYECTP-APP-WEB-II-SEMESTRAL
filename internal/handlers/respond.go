package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
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

// statusDeError convierte un error en el código HTTP apropiado.
// Usa coincidencias de texto porque este paquete no expone errores de servicio específicos.
func statusDeError(err error) int {
	if err == nil {
		return http.StatusInternalServerError
	}

	mensaje := strings.ToLower(err.Error())
	switch {
	case strings.Contains(mensaje, "credencial") || strings.Contains(mensaje, "contraseña") || strings.Contains(mensaje, "invalid credentials") || strings.Contains(mensaje, "incorrect password"):
		return http.StatusUnauthorized
	case strings.Contains(mensaje, "email") && strings.Contains(mensaje, "uso"):
		return http.StatusConflict
	case strings.Contains(mensaje, "no encontrada") || strings.Contains(mensaje, "no encontrado") || strings.Contains(mensaje, "not found") || strings.Contains(mensaje, "not exist"):
		return http.StatusNotFound
	case strings.Contains(mensaje, "requerido") || strings.Contains(mensaje, "invalido") || strings.Contains(mensaje, "invalid") || strings.Contains(mensaje, "mal formado") || strings.Contains(mensaje, "bad request") || strings.Contains(mensaje, "debe ser"):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
