package handlers

import (
	"github.com/go-chi/chi/v5"

	"condominio-api/internal/middleware"
)

// MontarRutasObligaciones registra las rutas de Obligaciones y Multas.
// Lectura (GET) abierta a cualquier usuario autenticado; escritura
// (POST/PUT/DELETE) restringida a rol "admin".
func MontarRutasObligaciones(r chi.Router, srv *Server) {
	soloAdmin := middleware.RequireRol("admin")

	// Obligaciones
	r.Get("/obligaciones", srv.ListarObligaciones)
	r.Get("/obligaciones/{id}", srv.GetObligacion)
	r.With(soloAdmin).Post("/obligaciones", srv.CreateObligacion)
	r.With(soloAdmin).Put("/obligaciones/{id}", srv.UpdateObligacion)
	r.With(soloAdmin).Delete("/obligaciones/{id}", srv.DeleteObligacion)

	// Multas
	r.Get("/multas", srv.ListarMultas)
	r.Get("/multas/{id}", srv.GetMulta)
	r.With(soloAdmin).Post("/multas", srv.CreateMulta)
	r.With(soloAdmin).Put("/multas/{id}", srv.UpdateMulta)
	r.With(soloAdmin).Delete("/multas/{id}", srv.DeleteMulta)
}
