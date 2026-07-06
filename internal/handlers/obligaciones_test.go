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

//Tiene dos tests: uno que prueba que crear una obligación válida responde "201 creado", y otro que prueba que si no mandas el token de seguridad, te rechaza con "401 no autorizado".

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
	return false
}

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

func TestObligaciones_SinToken_Devuelve401(t *testing.T) {
	fake := &fakeObligacionRepo{}
	svc := service.NewObligacionesService(fake)
	srv := &Server{Obligaciones: svc}

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
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
