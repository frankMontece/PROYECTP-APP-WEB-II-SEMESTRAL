package middleware

import (
	"condominio-api/internal/service"
	"context"
	"net/http"
	"strings"
)

type claveCtx string

const ClaveUsuarioID claveCtx = "usuarioID"

// Auth valida el token JWT y lo inyecta en el context
func Auth(auth *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			partes := strings.SplitN(header, " ", 2)

			if len(partes) != 2 || partes[0] != "Bearer" {
				responderNoAutorizado(w)
				return
			}

			userID, err := auth.ValidarToken(partes[1])
			if err != nil {
				responderNoAutorizado(w)
				return
			}

			ctx := context.WithValue(r.Context(), ClaveUsuarioID, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func responderNoAutorizado(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(`{"error": "token inexistente o inválido"}`))
}
