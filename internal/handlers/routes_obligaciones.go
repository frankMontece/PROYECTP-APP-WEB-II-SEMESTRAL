package handlers

import (
	"github.com/go-chi/chi/v5"
)

// MontarRutasObligaciones registra las rutas de Obligaciones y Multas.
// No monta auth ni decide middleware — eso lo hace main.go, igual que
// con MontarRutasSocial y MontarRutasParqueo.
func MontarRutasObligaciones(r chi.Router, srv *Server) {
	// Obligaciones
	r.Get("/obligaciones", srv.ListarObligaciones)
	r.Post("/obligaciones", srv.CreateObligacion)
	r.Get("/obligaciones/{id}", srv.GetObligacion)
	r.Put("/obligaciones/{id}", srv.UpdateObligacion)
	r.Delete("/obligaciones/{id}", srv.DeleteObligacion)

	// Multas
	r.Get("/multas", srv.ListarMultas)
	r.Post("/multas", srv.CreateMulta)
	r.Get("/multas/{id}", srv.GetMulta)
	r.Put("/multas/{id}", srv.UpdateMulta)
	r.Delete("/multas/{id}", srv.DeleteMulta)
}
