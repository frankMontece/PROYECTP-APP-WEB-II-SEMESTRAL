package handlers

import (
	"github.com/go-chi/chi/v5"
)

func (s *Server) MontarRutasParqueo(r chi.Router) {
	r.Get("/vehiculos", s.ListarVehiculos)
	r.Post("/vehiculos", s.CrearVehiculo)
	r.Get("/vehiculos/{id}", s.ObtenerVehiculo)
	r.Put("/vehiculos/{id}", s.ActualizarVehiculo)
	r.Delete("/vehiculos/{id}", s.BorrarVehiculo)

	r.Get("/visitas", s.ListarVisitas)
	r.Post("/visitas", s.CrearVisita)
	r.Get("/visitas/{id}", s.ObtenerVisita)
	r.Put("/visitas/{id}/entrada", s.ActualizarVisita)
	r.Put("/visitas/{id}/salida", s.RegistrarSalida)
	r.Delete("/visitas/{id}", s.BorrarVisita)

	r.Get("/accesos", s.ListarAccesos)
	r.Post("/accesos", s.CrearAcceso)
	r.Get("/accesos/{id}", s.ObtenerAcceso)
	r.Delete("/accesos/{id}", s.BorrarAcceso)
}
