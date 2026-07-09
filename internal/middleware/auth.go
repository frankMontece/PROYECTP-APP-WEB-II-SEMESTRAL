package middleware

import (
	"net/http"

	"condominio-api/internal/service"

	"context"
	"strings"
)

type claveContext string

const (
	ClaveUsuarioID claveContext = "usuarioID"
	ClaveRol       claveContext = "rol"
)

// Auth valida el JWT y coloca el ID de usuario y su rol en el contexto.
func Auth(auth *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			encabezado := r.Header.Get("Authorization")
			partes := strings.SplitN(encabezado, " ", 2)
			if len(partes) != 2 || partes[0] != "Bearer" {
				responderNoAutorizado(w)
				return
			}
			usuarioID, rol, err := auth.ValidarToken(partes[1])
			if err != nil {
				responderNoAutorizado(w)
				return
			}
			ctx := context.WithValue(r.Context(), ClaveUsuarioID, usuarioID)
			ctx = context.WithValue(ctx, ClaveRol, rol)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRol exige que el rol propagado por Auth esté dentro de los
// roles permitidos. Debe montarse DESPUÉS de Auth en la cadena de
// middlewares, porque depende de ClaveRol ya estar en el contexto.
//
// Uso:
//
//	r.Use(middleware.Auth(authService))
//	r.With(middleware.RequireRol("admin")).Delete("/vehiculos/{id}", srv.BorrarVehiculo)
func RequireRol(rolesPermitidos ...string) func(http.Handler) http.Handler {
	permitidos := make(map[string]bool, len(rolesPermitidos))
	for _, r := range rolesPermitidos {
		permitidos[r] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rol, ok := r.Context().Value(ClaveRol).(string)
			if !ok || !permitidos[rol] {
				responderProhibido(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func responderNoAutorizado(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(`{"error": "Token inexistente o invalido"}`))
}

func responderProhibido(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	_, _ = w.Write([]byte(`{"error": "No tienes permisos para realizar esta acción"}`))
}
