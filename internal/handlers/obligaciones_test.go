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
	for i, o := range f.obligaciones {
		if o.ID == id {
			datos.ID = id
			f.obligaciones[i] = datos
			return datos, true
		}
	}
	return models.Obligacion{}, false
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

// ---------- ListarObligaciones ----------

func TestListarObligaciones_DevuelveOK(t *testing.T) {
	fake := &fakeObligacionRepo{
		obligaciones: []models.Obligacion{{ID: 1, ResidenteID: 1, Monto: 50}},
	}
	svc := service.NewObligacionesService(fake)
	srv := &Server{Obligaciones: svc}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/obligaciones", nil)
	rec := httptest.NewRecorder()

	srv.ListarObligaciones(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("esperaba 200 OK, obtuve %d", rec.Code)
	}
}

// ---------- GetObligacion ----------

func TestGetObligacion_Encontrada_DevuelveOK(t *testing.T) {
	fake := &fakeObligacionRepo{
		obligaciones: []models.Obligacion{{ID: 1, ResidenteID: 1, Monto: 50}},
	}
	svc := service.NewObligacionesService(fake)
	srv := &Server{Obligaciones: svc}

	r := chi.NewRouter()
	r.Get("/api/v1/obligaciones/{id}", srv.GetObligacion)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/obligaciones/1", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("esperaba 200 OK, obtuve %d. Body: %s", rec.Code, rec.Body.String())
	}
}

func TestGetObligacion_NoEncontrada_Devuelve404(t *testing.T) {
	fake := &fakeObligacionRepo{}
	svc := service.NewObligacionesService(fake)
	srv := &Server{Obligaciones: svc}

	r := chi.NewRouter()
	r.Get("/api/v1/obligaciones/{id}", srv.GetObligacion)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/obligaciones/999", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("esperaba 404 Not Found, obtuve %d. Body: %s", rec.Code, rec.Body.String())
	}
}

func TestGetObligacion_IDInvalido_Devuelve400(t *testing.T) {
	fake := &fakeObligacionRepo{}
	svc := service.NewObligacionesService(fake)
	srv := &Server{Obligaciones: svc}

	r := chi.NewRouter()
	r.Get("/api/v1/obligaciones/{id}", srv.GetObligacion)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/obligaciones/abc", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("esperaba 400 Bad Request, obtuve %d", rec.Code)
	}
}

// ---------- CreateObligacion - casos de error ----------

func TestCreateObligacion_JSONInvalido_Devuelve400(t *testing.T) {
	fake := &fakeObligacionRepo{}
	svc := service.NewObligacionesService(fake)
	srv := &Server{Obligaciones: svc}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/obligaciones", bytes.NewReader([]byte("{esto no es json")))
	rec := httptest.NewRecorder()

	srv.CreateObligacion(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("esperaba 400 Bad Request, obtuve %d", rec.Code)
	}
}

func TestCreateObligacion_MontoInvalido_Devuelve400(t *testing.T) {
	fake := &fakeObligacionRepo{}
	svc := service.NewObligacionesService(fake)
	srv := &Server{Obligaciones: svc}

	body := models.Obligacion{
		ResidenteID: 1,
		Tipo:        "mensual",
		Monto:       -10.00, // inválido
		Periodo:     "2026-06",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/obligaciones", bytes.NewReader(jsonBody))
	rec := httptest.NewRecorder()

	srv.CreateObligacion(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("esperaba 400 Bad Request, obtuve %d. Body: %s", rec.Code, rec.Body.String())
	}
}

// ---------- UpdateObligacion ----------

func TestUpdateObligacion_Exitosa_DevuelveOK(t *testing.T) {
	fake := &fakeObligacionRepo{
		obligaciones: []models.Obligacion{{ID: 1, ResidenteID: 1, Tipo: "mensual", Monto: 50, Periodo: "2026-06"}},
	}
	svc := service.NewObligacionesService(fake)
	srv := &Server{Obligaciones: svc}

	r := chi.NewRouter()
	r.Put("/api/v1/obligaciones/{id}", srv.UpdateObligacion)

	body := models.Obligacion{ResidenteID: 1, Tipo: "mensual", Monto: 80, Periodo: "2026-06"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/obligaciones/1", bytes.NewReader(jsonBody))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("esperaba 200 OK, obtuve %d. Body: %s", rec.Code, rec.Body.String())
	}
}

func TestUpdateObligacion_IDInvalido_Devuelve400(t *testing.T) {
	fake := &fakeObligacionRepo{}
	svc := service.NewObligacionesService(fake)
	srv := &Server{Obligaciones: svc}

	r := chi.NewRouter()
	r.Put("/api/v1/obligaciones/{id}", srv.UpdateObligacion)

	body := models.Obligacion{ResidenteID: 1, Tipo: "mensual", Monto: 80, Periodo: "2026-06"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/obligaciones/abc", bytes.NewReader(jsonBody))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("esperaba 400 Bad Request, obtuve %d", rec.Code)
	}
}

// ---------- DeleteObligacion ----------

func TestDeleteObligacion_IDInvalido_Devuelve400(t *testing.T) {
	fake := &fakeObligacionRepo{}
	svc := service.NewObligacionesService(fake)
	srv := &Server{Obligaciones: svc}

	r := chi.NewRouter()
	r.Delete("/api/v1/obligaciones/{id}", srv.DeleteObligacion)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/obligaciones/abc", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("esperaba 400 Bad Request, obtuve %d", rec.Code)
	}
}
