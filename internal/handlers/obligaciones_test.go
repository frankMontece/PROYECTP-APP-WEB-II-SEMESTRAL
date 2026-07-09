package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"

	"condominio-api/internal/middleware"
	"condominio-api/internal/models"
	"condominio-api/internal/service"
	"condominio-api/internal/storage"
)

// hashPasswordParaTest genera un hash bcrypt real, para poder hacer
// AuthService.Login(...) con la misma contraseña en texto plano en tests.
func hashPasswordParaTest(t *testing.T, password string) (string, error) {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

type fakeObligacionRepo struct {
	obligaciones []models.Obligacion
}

var _ storage.ObligacionRepository = (*fakeObligacionRepo)(nil)

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

// nuevoServidorObligaciones construye un *Server solo con el servicio de
// Obligaciones montado, para los tests que llaman al handler directo
// (sin pasar por middleware).
func nuevoServidorObligaciones(repo storage.ObligacionRepository) *Server {
	svc := service.NewObligacionesService(repo)
	return NewServer(Services{
		Obligaciones: svc,
	})
}

// nuevoRouterObligacionesConRol monta el router completo (middleware.Auth +
// RequireRol reales) y devuelve un token JWT válido para el rol pedido.
// Reutiliza fakeUserRepo, ya definido en area_test.go.
func nuevoRouterObligacionesConRol(t *testing.T, rol string, repoObligaciones storage.ObligacionRepository) (chi.Router, string) {
	t.Helper()

	email := "obligaciones-" + rol + "@test.com"
	fakeUsuarios := &fakeUserRepo{
		usuarios: []models.Usuario{
			{ID: 1, Email: email, Rol: rol},
		},
	}
	// Generamos un hash real para poder hacer Login con el mismo password.
	hash, err := hashPasswordParaTest(t, "password123")
	if err != nil {
		t.Fatalf("error generando hash de test: %v", err)
	}
	fakeUsuarios.usuarios[0].PasswordHash = hash

	authSvc := service.NewAuthService(fakeUsuarios)
	svc := service.NewObligacionesService(repoObligaciones)
	srv := NewServer(Services{
		Auth:         authSvc,
		Obligaciones: svc,
	})

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authSvc))
			MontarRutasObligaciones(r, srv)
		})
	})

	token, err := authSvc.Login(email, "password123")
	if err != nil {
		t.Fatalf("error generando token de test: %v", err)
	}

	return r, token
}

// ── Tests ─────────────────────────────────────────────────────────────────────

func TestCreateObligacion_DevuelveCreated(t *testing.T) {
	fake := &fakeObligacionRepo{}
	srv := nuevoServidorObligaciones(fake)

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

// Test de autenticación real: usa el middleware.Auth de verdad.
func TestObligaciones_SinToken_Devuelve401(t *testing.T) {
	fake := &fakeObligacionRepo{}
	r, _ := nuevoRouterObligacionesConRol(t, "residente", fake)

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
	srv := nuevoServidorObligaciones(fake)

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
	srv := nuevoServidorObligaciones(fake)

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
	srv := nuevoServidorObligaciones(fake)

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
	srv := nuevoServidorObligaciones(fake)

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
	srv := nuevoServidorObligaciones(fake)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/obligaciones", bytes.NewReader([]byte("{esto no es json")))
	rec := httptest.NewRecorder()

	srv.CreateObligacion(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("esperaba 400 Bad Request, obtuve %d", rec.Code)
	}
}

func TestCreateObligacion_MontoInvalido_Devuelve400(t *testing.T) {
	fake := &fakeObligacionRepo{}
	srv := nuevoServidorObligaciones(fake)

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
	srv := nuevoServidorObligaciones(fake)

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
	srv := nuevoServidorObligaciones(fake)

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
	srv := nuevoServidorObligaciones(fake)

	r := chi.NewRouter()
	r.Delete("/api/v1/obligaciones/{id}", srv.DeleteObligacion)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/obligaciones/abc", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("esperaba 400 Bad Request, obtuve %d", rec.Code)
	}
}

// ---------- RequireRol — protección por rol en escritura, con JWT real ────────

func TestCreateObligacion_ConTokenResidente_Devuelve403(t *testing.T) {
	fake := &fakeObligacionRepo{}
	r, token := nuevoRouterObligacionesConRol(t, service.RolResidente, fake)

	body := models.Obligacion{
		ResidenteID: 1,
		Tipo:        "mensual",
		Monto:       150.00,
		Periodo:     "2026-06",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/obligaciones", bytes.NewReader(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("esperaba 403 Forbidden con rol residente, obtuve %d. Body: %s", rec.Code, rec.Body.String())
	}
}

func TestCreateObligacion_ConTokenAdmin_Devuelve201(t *testing.T) {
	fake := &fakeObligacionRepo{}
	r, token := nuevoRouterObligacionesConRol(t, service.RolAdmin, fake)

	body := models.Obligacion{
		ResidenteID: 1,
		Tipo:        "mensual",
		Monto:       150.00,
		Periodo:     "2026-06",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/obligaciones", bytes.NewReader(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("esperaba 201 Created con rol admin, obtuve %d. Body: %s", rec.Code, rec.Body.String())
	}
}

func TestListarObligaciones_ConTokenResidente_Devuelve200(t *testing.T) {
	fake := &fakeObligacionRepo{
		obligaciones: []models.Obligacion{{ID: 1, ResidenteID: 1, Monto: 50}},
	}
	r, token := nuevoRouterObligacionesConRol(t, service.RolResidente, fake)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/obligaciones", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("esperaba 200 OK para lectura con rol residente, obtuve %d", rec.Code)
	}
}
