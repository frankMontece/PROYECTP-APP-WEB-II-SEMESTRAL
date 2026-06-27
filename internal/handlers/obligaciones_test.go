package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"condominio-api/internal/models"
	"condominio-api/internal/service"
)

// fakeObligacionRepo SÍ guarda datos (a diferencia del mock del Test 1),
// pero en un slice en memoria — nunca toca SQLite real. Por eso es un
// "fake", no un mock: simula el comportamiento completo, no solo verifica una llamada.
type fakeObligacionRepo struct {
	obligaciones []models.Obligacion
}

func (f *fakeObligacionRepo) ListarObligaciones() []models.Obligacion {
	return f.obligaciones
}

func (f *fakeObligacionRepo) BuscarObligacionPorID(id uint) (models.Obligacion, bool) {
	for _, o := range f.obligaciones {
		if o.ID == id {
			return o, true
		}
	}
	return models.Obligacion{}, false
}

func (f *fakeObligacionRepo) CrearObligacion(o models.Obligacion) models.Obligacion {
	o.ID = uint(len(f.obligaciones) + 1)
	f.obligaciones = append(f.obligaciones, o)
	return o
}

func (f *fakeObligacionRepo) ActualizarObligacion(id uint, datos models.Obligacion) (models.Obligacion, bool) {
	return models.Obligacion{}, false // no usado en estos tests
}

func (f *fakeObligacionRepo) BorrarObligacion(id uint) bool {
	return false // no usado en estos tests
}

// TestCreateObligacion_DevuelveCreated prueba el camino feliz del handler:
// con datos válidos, debe responder 201 Created.
// Llama directo al handler (sin pasar por router ni middleware) — es un
// test aislado de la función CreateObligacion únicamente.
func TestCreateObligacion_DevuelveCreated(t *testing.T) {
	fake := &fakeObligacionRepo{}
	svc := service.NewObligacionesService(fake)
	srv := &Server{Obligaciones: svc}

	body := models.Obligacion{
		ResidenteID: 1,
		Tipo:        "mensual",
		Monto:       150.00,
		Periodo:     "2026-06",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/obligaciones", bytes.NewReader(jsonBody))
	rec := httptest.NewRecorder()

	srv.CreateObligacion(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("esperaba 201 Created, obtuve %d. Body: %s", rec.Code, rec.Body.String())
	}
}

// TestObligaciones_SinToken_Devuelve401 prueba que una ruta protegida
// SIN el header Authorization responde 401. A diferencia del test anterior,
// aquí se monta un router chi real con un middleware simulado, porque
// el 401 no lo genera el handler — lo genera el middleware ANTES de llegar
// al handler. Por eso este test pasa por el router completo.
func TestObligaciones_SinToken_Devuelve401(t *testing.T) {
	fake := &fakeObligacionRepo{}
	svc := service.NewObligacionesService(fake)
	srv := &Server{Obligaciones: svc}

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			// Simulación del middleware real (internal/middelware/auth.go):
			// la misma lógica — sin header Authorization, corta con 401
			// y nunca llama a next.ServeHTTP.
			r.Use(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					if req.Header.Get("Authorization") == "" {
						w.WriteHeader(http.StatusUnauthorized)
						return
					}
					next.ServeHTTP(w, req)
				})
			})
			r.Get("/obligaciones", srv.ListarObligaciones)
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/obligaciones", nil)
	// A propósito NO se agrega el header Authorization
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("esperaba 401 Unauthorized sin token, obtuve %d", rec.Code)
	}
}
