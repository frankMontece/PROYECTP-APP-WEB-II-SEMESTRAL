package handlers_test

// Patrón: handler con httptest + fake en memoria

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"condominio-api/internal/handlers"
	"condominio-api/internal/middleware"
	"condominio-api/internal/models"
	"condominio-api/internal/service"
)

// ── Fake del repositorio de usuarios ──

type fakeUserRepo struct{}

func (f *fakeUserRepo) CrearUsuario(u models.Usuario) (models.Usuario, error) {
	u.ID = 1
	return u, nil
}

func (f *fakeUserRepo) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	return models.Usuario{}, false
}

// ── Fake del almacén de parqueo ──

type fakeAlmacenParqueo struct {
	vehiculos []models.Vehiculo
	nextID    uint
}

// Vehiculos
func (f *fakeAlmacenParqueo) ListarVehiculos() []models.Vehiculo { return f.vehiculos }

func (f *fakeAlmacenParqueo) BuscarVehiculoPorID(id uint) (models.Vehiculo, bool) {
	for _, v := range f.vehiculos {
		if v.ID == id {
			return v, true
		}
	}
	return models.Vehiculo{}, false
}

func (f *fakeAlmacenParqueo) CrearVehiculo(v models.Vehiculo) models.Vehiculo {
	f.nextID++
	v.ID = f.nextID
	f.vehiculos = append(f.vehiculos, v)
	return v
}

func (f *fakeAlmacenParqueo) ActualizarVehiculo(id uint, datos models.Vehiculo) (models.Vehiculo, bool) {
	for i, v := range f.vehiculos {
		if v.ID == id {
			datos.ID = id
			f.vehiculos[i] = datos
			return datos, true
		}
	}
	return models.Vehiculo{}, false
}

func (f *fakeAlmacenParqueo) BorrarVehiculo(id uint) bool {
	for i, v := range f.vehiculos {
		if v.ID == id {
			f.vehiculos = append(f.vehiculos[:i], f.vehiculos[i+1:]...)
			return true
		}
	}
	return false
}

// Visitas (no se prueban aquí, implementación mínima para satisfacer la interfaz)
func (f *fakeAlmacenParqueo) ListarVisitas() []models.VisitaVehiculo { return nil }
func (f *fakeAlmacenParqueo) BuscarVisitaPorID(id uint) (models.VisitaVehiculo, bool) {
	return models.VisitaVehiculo{}, false
}
func (f *fakeAlmacenParqueo) CrearVisita(v models.VisitaVehiculo) models.VisitaVehiculo { return v }
func (f *fakeAlmacenParqueo) ActualizarVisita(id uint, datos models.VisitaVehiculo) (models.VisitaVehiculo, bool) {
	return models.VisitaVehiculo{}, false
}
func (f *fakeAlmacenParqueo) RegistrarSalidaVisita(id uint) (models.VisitaVehiculo, bool) {
	return models.VisitaVehiculo{}, false
}
func (f *fakeAlmacenParqueo) BorrarVisita(id uint) bool { return false }

// Accesos (no se prueban aquí, implementación mínima para satisfacer la interfaz)
func (f *fakeAlmacenParqueo) ListarAccesos() []models.AccesoVehiculo { return nil }
func (f *fakeAlmacenParqueo) BuscarAccesoPorID(id uint) (models.AccesoVehiculo, bool) {
	return models.AccesoVehiculo{}, false
}
func (f *fakeAlmacenParqueo) CrearAcceso(a models.AccesoVehiculo) models.AccesoVehiculo { return a }
func (f *fakeAlmacenParqueo) BorrarAcceso(id uint) bool                                 { return false }

// ── Helper: construye router + devuelve token de prueba ──

func setupRouter(t *testing.T) (*chi.Mux, string) {
	t.Helper()

	fake := &fakeAlmacenParqueo{}
	authService := service.NewAuthService(&fakeUserRepo{}, []byte("test-secret"), 2*time.Hour)

	servidor := handlers.NewServer(
		handlers.Services{
			Auth:      authService,
			Vehiculos: service.NewVehiculoService(fake),
			Visitas:   service.NewVisitaService(fake),
			Accesos:   service.NewAccesoService(fake),
		},
	)

	// Token válido de prueba: usuario con ID 1, sin contraseña real
	token, err := authService.GenerarToken(models.Usuario{ID: 1})
	require.NoError(t, err)

	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(authService))
		r.Get("/vehiculos", servidor.ListarVehiculos)
		r.Post("/vehiculos", servidor.CrearVehiculo)
		r.Get("/vehiculos/{id}", servidor.ObtenerVehiculo)
		r.Put("/vehiculos/{id}", servidor.ActualizarVehiculo)
		r.Delete("/vehiculos/{id}", servidor.BorrarVehiculo)
	})

	return r, token
}

// ── Test 1: POST /vehiculos con token → 201 ───────────────────────────────────

func TestCrearVehiculo_ConToken_201(t *testing.T) {
	// Preparar
	r, token := setupRouter(t)

	body := models.Vehiculo{
		ResidenteID: 1,
		Placa:       "PBG-2241",
		Marca:       "Toyota",
		Modelo:      "Corolla",
		Color:       "Blanco",
	}
	buf, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/vehiculos", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	// Ejecutar
	r.ServeHTTP(rec, req)

	// Verificar
	require.Equal(t, http.StatusCreated, rec.Code)

	var respuesta models.Vehiculo
	err := json.Unmarshal(rec.Body.Bytes(), &respuesta)
	require.NoError(t, err)
	assert.Equal(t, "PBG-2241", respuesta.Placa)
	assert.Equal(t, "Toyota", respuesta.Marca)
	assert.True(t, respuesta.Activo) // el service asigna Activo=true
	assert.NotZero(t, respuesta.ID)  // el fake asigna un ID
}

// ── Test 2: POST /vehiculos SIN token → 401 ──

func TestCrearVehiculo_SinToken_401(t *testing.T) {
	// Preparar
	r, _ := setupRouter(t) // el token no se usa en este test

	body := models.Vehiculo{ResidenteID: 1, Placa: "PBG-2241", Marca: "Toyota"}
	buf, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/vehiculos", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	// Sin "Authorization" header: el middleware debe rechazar
	rec := httptest.NewRecorder()

	// Ejecutar
	r.ServeHTTP(rec, req)

	// Verificar
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
