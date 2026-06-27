package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// parseID extrae el parámetro "id" de la URL y lo convierte a uint.
func parseID(r *http.Request) (uint, error) {
	id64, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id64), nil
}
