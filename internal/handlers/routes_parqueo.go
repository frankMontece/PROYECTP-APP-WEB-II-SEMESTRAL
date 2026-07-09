package handlers

import (
	"github.com/go-chi/chi/v5"

	"condominio-api/internal/middleware"
)

// MontarRutasParqueo registra las rutas de Vehículos, Visitas y Accesos.
// Lectura (GET) abierta a cualquier usuario autenticado; escritura
// (POST/PUT/DELETE) restringida a rol "admin".
func MontarRutasParqueo(r chi.Router, s *Server) {
	soloAdmin := middleware.RequireRol("admin")

	// Vehículos
	r.Get("/vehiculos", s.ListarVehiculos)
	r.Get("/vehiculos/{id}", s.ObtenerVehiculo)
	r.With(soloAdmin).Post("/vehiculos", s.CrearVehiculo)
	r.With(soloAdmin).Put("/vehiculos/{id}", s.ActualizarVehiculo)
	r.With(soloAdmin).Delete("/vehiculos/{id}", s.BorrarVehiculo)

	// Visitas
	r.Get("/visitas", s.ListarVisitas)
	r.Get("/visitas/{id}", s.ObtenerVisita)
	r.With(soloAdmin).Post("/visitas", s.CrearVisita)
	r.With(soloAdmin).Put("/visitas/{id}/entrada", s.ActualizarVisita)
	r.With(soloAdmin).Put("/visitas/{id}/salida", s.RegistrarSalida)
	r.With(soloAdmin).Delete("/visitas/{id}", s.BorrarVisita)

	// Accesos
	r.Get("/accesos", s.ListarAccesos)
	r.Get("/accesos/{id}", s.ObtenerAcceso)
	r.With(soloAdmin).Post("/accesos", s.CrearAcceso)
	r.With(soloAdmin).Delete("/accesos/{id}", s.BorrarAcceso)
}
