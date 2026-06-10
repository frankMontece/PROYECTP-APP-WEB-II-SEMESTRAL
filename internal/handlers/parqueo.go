package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"condominio-api/internal/models"
	"condominio-api/internal/storage"
)

func parseID(r *http.Request) (uint, bool) {
	idStr := chi.URLParam(r, "id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, false
	}
	return uint(id64), true
}

func MontarRutasParqueo(r chi.Router, store storage.AlmacenParqueo) {
 
	// --- Vehículos de residentes ---
	r.Get("/vehiculos", func(w http.ResponseWriter, req *http.Request) {
		GetAllVehiculos(w, req, store)
	})
	r.Post("/vehiculos", func(w http.ResponseWriter, req *http.Request) {
		CreateVehiculo(w, req, store)
	})
	r.Get("/vehiculos/{id}", func(w http.ResponseWriter, req *http.Request) {
		GetVehiculo(w, req, store)
	})
	r.Put("/vehiculos/{id}", func(w http.ResponseWriter, req *http.Request) {
		UpdateVehiculo(w, req, store)
	})
	r.Delete("/vehiculos/{id}", func(w http.ResponseWriter, req *http.Request) {
		DeleteVehiculo(w, req, store)
	})
 
	// --- Visitas de vehículos externos ---
	r.Get("/visitas", func(w http.ResponseWriter, req *http.Request) {
		GetAllVisitas(w, req, store)
	})
	r.Post("/visitas", func(w http.ResponseWriter, req *http.Request) {
		CreateVisita(w, req, store)
	})
	r.Get("/visitas/{id}", func(w http.ResponseWriter, req *http.Request) {
		GetVisita(w, req, store)
	})
	r.Put("/visitas/{id}", func(w http.ResponseWriter, req *http.Request) {
		UpdateVisita(w, req, store)
	})
	r.Post("/visitas/{id}/salida", func(w http.ResponseWriter, req *http.Request) {
		RegistrarSalida(w, req, store)
	})
	r.Delete("/visitas/{id}", func(w http.ResponseWriter, req *http.Request) {
		DeleteVisita(w, req, store)
	})
 
	// --- Bitácora de accesos de vehículos de residentes ---
	r.Get("/accesos", func(w http.ResponseWriter, req *http.Request) {
		GetAllAccesos(w, req, store)
	})
	r.Post("/accesos", func(w http.ResponseWriter, req *http.Request) {
		CreateAcceso(w, req, store)
	})
	r.Get("/accesos/{id}", func(w http.ResponseWriter, req *http.Request) {
		GetAcceso(w, req, store)
	})
	r.Delete("/accesos/{id}", func(w http.ResponseWriter, req *http.Request) {
		DeleteAcceso(w, req, store)
	})
}