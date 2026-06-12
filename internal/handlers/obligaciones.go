package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"condominio-api/internal/models"
	"condominio-api/internal/storage"
)

// ROUTER, basicamente asigna cada endpoint a su función handler correspondiente, pasando el store para acceso a datos.

func MontarRutasAlicuotas(r chi.Router, store storage.AlmacenObligaciones) {

	// Obligaciones
	r.Get("/obligaciones", func(w http.ResponseWriter, req *http.Request) {
		GetAllObligaciones(w, req, store)
	})
	r.Post("/obligaciones", func(w http.ResponseWriter, req *http.Request) {
		CreateObligacion(w, req, store)
	})
	r.Get("/obligaciones/{id}", func(w http.ResponseWriter, req *http.Request) {
		GetObligacion(w, req, store)
	})
	r.Put("/obligaciones/{id}", func(w http.ResponseWriter, req *http.Request) {
		UpdateObligacion(w, req, store)
	})
	r.Delete("/obligaciones/{id}", func(w http.ResponseWriter, req *http.Request) {
		DeleteObligacion(w, req, store)
	})

	// Multas
	r.Get("/multas", func(w http.ResponseWriter, req *http.Request) {
		GetAllMultas(w, req, store)
	})
	r.Post("/multas", func(w http.ResponseWriter, req *http.Request) {
		CreateMulta(w, req, store)
	})
	r.Get("/multas/{id}", func(w http.ResponseWriter, req *http.Request) {
		GetMulta(w, req, store)
	})
	r.Put("/multas/{id}", func(w http.ResponseWriter, req *http.Request) {
		UpdateMulta(w, req, store)
	})
	r.Delete("/multas/{id}", func(w http.ResponseWriter, req *http.Request) {
		DeleteMulta(w, req, store)
	})
}

// parseIDAlicuotas lee el parámetro {id} de la URL y lo convierte a uint.
// Devuelve (id, true) si es válido o (0, false) si no lo es.
func parseIDAlicuotas(r *http.Request) (uint, bool) {
	idStr := chi.URLParam(r, "id")
	var id uint
	_, err := fmt.Sscanf(idStr, "%d", &id)
	if err != nil || id == 0 {
		return 0, false
	}
	return id, true
}

// OBLIGACIONES, cada función atiende un endpoint específico, valida la entrada, llama al store y responde con JSON o error según corresponda.

// GetAllObligaciones atiende GET /api/v1/obligaciones
func GetAllObligaciones(w http.ResponseWriter, r *http.Request, store storage.AlmacenObligaciones) {
	lista, err := store.ListarObligaciones()
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, lista)
}

// GetObligacion atiende GET /api/v1/obligaciones/{id}
func GetObligacion(w http.ResponseWriter, r *http.Request, store storage.AlmacenObligaciones) {
	id, ok := parseIDAlicuotas(r)
	if !ok {
		RespondError(w, http.StatusBadRequest, "ID debe ser un número entero positivo")
		return
	}

	obligacion, err := store.BuscarObligacionPorID(id)
	if err != nil {
		RespondError(w, http.StatusNotFound, err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, obligacion)
}

// CreateObligacion atiende POST /api/v1/obligaciones
func CreateObligacion(w http.ResponseWriter, r *http.Request, store storage.AlmacenObligaciones) {
	var body models.Obligacion
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}

	if body.ResidenteID == 0 {
		RespondError(w, http.StatusBadRequest, "residente_id es requerido")
		return
	}
	if body.Monto <= 0 {
		RespondError(w, http.StatusBadRequest, "monto debe ser mayor a 0")
		return
	}
	if strings.TrimSpace(body.Periodo) == "" {
		RespondError(w, http.StatusBadRequest, "periodo es requerido (formato YYYY-MM)")
		return
	}
	if strings.TrimSpace(body.Tipo) == "" {
		RespondError(w, http.StatusBadRequest, "tipo es requerido")
		return
	}

	// El servidor asigna la fecha de emisión, no el cliente
	body.FechaEmision = time.Now()
	body.FechaPago = nil // null hasta que se registre el pago

	resultado, err := store.CrearObligacion(body)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, resultado)
}

// UpdateObligacion atiende PUT /api/v1/obligaciones/{id}

func UpdateObligacion(w http.ResponseWriter, r *http.Request, store storage.AlmacenObligaciones) {
	id, ok := parseIDAlicuotas(r)
	if !ok {
		RespondError(w, http.StatusBadRequest, "ID debe ser un número entero positivo")
		return
	}

	var body models.Obligacion
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}

	if body.ResidenteID == 0 {
		RespondError(w, http.StatusBadRequest, "residente_id es requerido")
		return
	}
	if body.Monto <= 0 {
		RespondError(w, http.StatusBadRequest, "monto debe ser mayor a 0")
		return
	}
	if strings.TrimSpace(body.Periodo) == "" {
		RespondError(w, http.StatusBadRequest, "periodo es requerido (formato YYYY-MM)")
		return
	}
	if strings.TrimSpace(body.Tipo) == "" {
		RespondError(w, http.StatusBadRequest, "tipo es requerido")
		return
	}

	resultado, err := store.ActualizarObligacion(id, body)
	if err != nil {
		RespondError(w, http.StatusNotFound, err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, resultado)
}

// DeleteObligacion atiende DELETE /api/v1/obligaciones/{id}
func DeleteObligacion(w http.ResponseWriter, r *http.Request, store storage.AlmacenObligaciones) {
	id, ok := parseIDAlicuotas(r)
	if !ok {
		RespondError(w, http.StatusBadRequest, "ID debe ser un número entero positivo")
		return
	}

	if err := store.EliminarObligacion(id); err != nil {
		RespondError(w, http.StatusNotFound, err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, map[string]string{"mensaje": "obligacion eliminada"})
}

// MULTAS, cada función atiende un endpoint específico, valida la entrada, llama al store y responde con JSON o error según corresponda.

// GetAllMultas atiende GET /api/v1/multas
func GetAllMultas(w http.ResponseWriter, r *http.Request, store storage.AlmacenObligaciones) {
	lista, err := store.ListarMultas()
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, lista)
}

// GetMulta atiende GET /api/v1/multas/{id}
func GetMulta(w http.ResponseWriter, r *http.Request, store storage.AlmacenObligaciones) {
	id, ok := parseIDAlicuotas(r)
	if !ok {
		RespondError(w, http.StatusBadRequest, "ID debe ser un número entero positivo")
		return
	}

	multa, err := store.BuscarMultaPorID(id)
	if err != nil {
		RespondError(w, http.StatusNotFound, err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, multa)
}

// CreateMulta atiende POST /api/v1/multas
func CreateMulta(w http.ResponseWriter, r *http.Request, store storage.AlmacenObligaciones) {
	var body models.Multa
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}

	if body.ResidenteID == 0 {
		RespondError(w, http.StatusBadRequest, "residente_id es requerido")
		return
	}
	if body.Monto <= 0 {
		RespondError(w, http.StatusBadRequest, "monto debe ser mayor a 0")
		return
	}
	if strings.TrimSpace(body.Motivo) == "" {
		RespondError(w, http.StatusBadRequest, "motivo es requerido")
		return
	}

	// El servidor asigna la fecha de emisión
	body.FechaEmision = time.Now()
	body.FechaPago = nil

	resultado, err := store.CrearMulta(body)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, resultado)
}

// UpdateMulta atiende PUT /api/v1/multas/{id}
func UpdateMulta(w http.ResponseWriter, r *http.Request, store storage.AlmacenObligaciones) {
	id, ok := parseIDAlicuotas(r)
	if !ok {
		RespondError(w, http.StatusBadRequest, "ID debe ser un número entero positivo")
		return
	}

	var body models.Multa
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}

	if body.ResidenteID == 0 {
		RespondError(w, http.StatusBadRequest, "residente_id es requerido")
		return
	}
	if body.Monto <= 0 {
		RespondError(w, http.StatusBadRequest, "monto debe ser mayor a 0")
		return
	}
	if strings.TrimSpace(body.Motivo) == "" {
		RespondError(w, http.StatusBadRequest, "motivo es requerido")
		return
	}

	resultado, err := store.ActualizarMulta(id, body)
	if err != nil {
		RespondError(w, http.StatusNotFound, err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, resultado)
}

// DeleteMulta atiende DELETE /api/v1/multas/{id}
func DeleteMulta(w http.ResponseWriter, r *http.Request, store storage.AlmacenObligaciones) {
	id, ok := parseIDAlicuotas(r)
	if !ok {
		RespondError(w, http.StatusBadRequest, "ID debe ser un número entero positivo")
		return
	}

	if err := store.EliminarMulta(id); err != nil {
		RespondError(w, http.StatusNotFound, err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, map[string]string{"mensaje": "multa eliminada"})
}
