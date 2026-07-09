package handlers_test

// Patrón: mismo enfoque que handlers_vehiculo_test.go (httptest + fake en
// memoria), pero autocontenido para el módulo Área Social:
//   - Reutiliza el tipo fakeUserRepo ya declarado en handlers_vehiculo_test.go
//     (mismo paquete handlers_test) — NO se vuelve a declarar aquí.
//   - No depende de setupRouterAdmin (no existe aún en el paquete). Se
//     definen setupRouterSocial y setupRouterSocialAdmin, con nombres únicos
//     para evitar colisión futura si se agrega esa función en otro archivo.

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

// ── Fake del almacén de Área Social (implementa las 3 interfaces de
//    storage/almacen-social.go: AreaSocialRepository, ReservaAreaRepository,
//    NotificacionRepository) ───────────────────────────────────────────────

type fakeAlmacenSocial struct {
	areas          []models.AreaSocial
	nextAreaID     uint
	reservas       []models.ReservaArea
	nextReservaID  uint
	notificaciones []models.Notificacion
	nextNotifID    uint
}

// Áreas Sociales
func (f *fakeAlmacenSocial) ListarAreas() []models.AreaSocial { return f.areas }

func (f *fakeAlmacenSocial) BuscarAreaPorID(id uint) (models.AreaSocial, bool) {
	for _, a := range f.areas {
		if a.ID == id {
			return a, true
		}
	}
	return models.AreaSocial{}, false
}

func (f *fakeAlmacenSocial) CrearArea(a models.AreaSocial) models.AreaSocial {
	f.nextAreaID++
	a.ID = f.nextAreaID
	f.areas = append(f.areas, a)
	return a
}

func (f *fakeAlmacenSocial) ActualizarArea(id uint, datos models.AreaSocial) (models.AreaSocial, bool) {
	for i, a := range f.areas {
		if a.ID == id {
			datos.ID = id
			f.areas[i] = datos
			return datos, true
		}
	}
	return models.AreaSocial{}, false
}

func (f *fakeAlmacenSocial) BorrarArea(id uint) bool {
	for i, a := range f.areas {
		if a.ID == id {
			f.areas = append(f.areas[:i], f.areas[i+1:]...)
			return true
		}
	}
	return false
}

// Reservas
func (f *fakeAlmacenSocial) ListarReservas() []models.ReservaArea { return f.reservas }

func (f *fakeAlmacenSocial) BuscarReservaPorID(id uint) (models.ReservaArea, bool) {
	for _, r := range f.reservas {
		if r.ID == id {
			return r, true
		}
	}
	return models.ReservaArea{}, false
}

func (f *fakeAlmacenSocial) CrearReserva(r models.ReservaArea) models.ReservaArea {
	f.nextReservaID++
	r.ID = f.nextReservaID
	f.reservas = append(f.reservas, r)
	return r
}

func (f *fakeAlmacenSocial) ActualizarReserva(id uint, datos models.ReservaArea) (models.ReservaArea, bool) {
	for i, r := range f.reservas {
		if r.ID == id {
			datos.ID = id
			f.reservas[i] = datos
			return datos, true
		}
	}
	return models.ReservaArea{}, false
}

func (f *fakeAlmacenSocial) BorrarReserva(id uint) bool {
	for i, r := range f.reservas {
		if r.ID == id {
			f.reservas = append(f.reservas[:i], f.reservas[i+1:]...)
			return true
		}
	}
	return false
}

// Notificaciones
func (f *fakeAlmacenSocial) ListarNotificaciones() []models.Notificacion { return f.notificaciones }

func (f *fakeAlmacenSocial) BuscarNotificacionPorID(id uint) (models.Notificacion, bool) {
	for _, n := range f.notificaciones {
		if n.ID == id {
			return n, true
		}
	}
	return models.Notificacion{}, false
}

func (f *fakeAlmacenSocial) CrearNotificacion(n models.Notificacion) models.Notificacion {
	f.nextNotifID++
	n.ID = f.nextNotifID
	f.notificaciones = append(f.notificaciones, n)
	return n
}

func (f *fakeAlmacenSocial) MarcarComoLeida(id uint) (models.Notificacion, bool) {
	for i, n := range f.notificaciones {
		if n.ID == id {
			f.notificaciones[i].Leida = true
			return f.notificaciones[i], true
		}
	}
	return models.Notificacion{}, false
}

func (f *fakeAlmacenSocial) BorrarNotificacion(id uint) bool {
	for i, n := range f.notificaciones {
		if n.ID == id {
			f.notificaciones = append(f.notificaciones[:i], f.notificaciones[i+1:]...)
			return true
		}
	}
	return false
}

// ── Helpers: construyen router + devuelven token de prueba ──────────────────
//
// Nota sobre nombres: fakeUserRepo NO se redeclara aquí; ya existe en
// handlers_vehiculo_test.go dentro del mismo paquete handlers_test.

// setupRouterSocial monta el router del módulo Área Social con un usuario
// autenticado normal (rol "residente"). Sirve para probar las rutas de
// lectura (GET) y la creación de reservas, abiertas a cualquier autenticado.
func setupRouterSocial(t *testing.T) (chi.Router, string) {
	t.Helper()

	repoUsuarios := &fakeUserRepo{}
	authService := service.NewAuthService(repoUsuarios)

	_, err := authService.Registrar("residente@test.com", "password123")
	require.NoError(t, err)

	token, err := authService.Login("residente@test.com", "password123")
	require.NoError(t, err)

	almacen := &fakeAlmacenSocial{}
	areaSvc := service.NewAreaService(almacen)
	reservaSvc := service.NewReservaService(almacen)
	notifSvc := service.NewNotificacionService(almacen)

	srv := handlers.NewServer(handlers.Services{
		Auth:         authService,
		Area:         areaSvc,
		Reserva:      reservaSvc,
		Notificacion: notifSvc,
	})

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authService))
			handlers.MontarRutasSocial(r, srv)
		})
	})

	return r, token
}

// setupRouterSocialAdmin monta el mismo router, pero con un usuario creado
// directamente con Rol: service.RolAdmin. Registrar() siempre asigna
// RolResidente por diseño (ver auth.go), así que el usuario admin se crea
// vía repo.CrearUsuario directamente, igual que lo haría un seed real.
func setupRouterSocialAdmin(t *testing.T) (chi.Router, string) {
	t.Helper()

	repoUsuarios := &fakeUserRepo{}
	authService := service.NewAuthService(repoUsuarios)

	_, err := repoUsuarios.CrearUsuario(models.Usuario{
		Email:        "admin@test.com",
		PasswordHash: hashPassword(t, "password123"),
		Rol:          service.RolAdmin,
	})
	require.NoError(t, err)

	token, err := authService.Login("admin@test.com", "password123")
	require.NoError(t, err)

	almacen := &fakeAlmacenSocial{}
	areaSvc := service.NewAreaService(almacen)
	reservaSvc := service.NewReservaService(almacen)
	notifSvc := service.NewNotificacionService(almacen)

	srv := handlers.NewServer(handlers.Services{
		Auth:         authService,
		Area:         areaSvc,
		Reserva:      reservaSvc,
		Notificacion: notifSvc,
	})

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authService))
			handlers.MontarRutasSocial(r, srv)
		})
	})

	return r, token
}

// hashPassword genera un hash bcrypt válido reutilizando el propio
// AuthService.Registrar (evita importar bcrypt de nuevo en el test),
// necesario porque CrearUsuario espera un PasswordHash ya calculado y
// Registrar normalmente asignaría el rol "residente".
func hashPassword(t *testing.T, plain string) string {
	t.Helper()
	tmpRepo := &fakeUserRepo{}
	tmpAuth := service.NewAuthService(tmpRepo)
	u, err := tmpAuth.Registrar("tmp-hash-helper@test.com", plain)
	require.NoError(t, err)
	return u.PasswordHash
}

// ═══════════════════════════════════════════════════════════════════════
// ÁREAS SOCIALES
// ═══════════════════════════════════════════════════════════════════════

func TestCrearArea_ConTokenAdmin_201(t *testing.T) {
	r, token := setupRouterSocialAdmin(t)

	body := models.AreaSocial{Nombre: "Salón Comunal", Descripcion: "Salón para eventos", Capacidad: 50}
	buf, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/areas-sociales", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	var respuesta models.AreaSocial
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &respuesta))
	assert.Equal(t, "Salón Comunal", respuesta.Nombre)
	assert.True(t, respuesta.Activo) // el service asigna Activo=true
	assert.NotZero(t, respuesta.ID)
}

func TestCrearArea_SinToken_401(t *testing.T) {
	r, _ := setupRouterSocial(t)

	body := models.AreaSocial{Nombre: "Salón Comunal", Capacidad: 50}
	buf, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/areas-sociales", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestCrearArea_TokenNoAdmin_403(t *testing.T) {
	// Rol "residente" no puede crear áreas; ruta protegida con RequireRol("admin")
	r, token := setupRouterSocial(t)

	body := models.AreaSocial{Nombre: "Salón Comunal", Capacidad: 50}
	buf, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/areas-sociales", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestCrearArea_NombreVacio_400(t *testing.T) {
	r, token := setupRouterSocialAdmin(t)

	body := models.AreaSocial{Nombre: "", Capacidad: 50}
	buf, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/areas-sociales", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCrearArea_CapacidadInvalida_400(t *testing.T) {
	r, token := setupRouterSocialAdmin(t)

	body := models.AreaSocial{Nombre: "Salón Comunal", Capacidad: 0}
	buf, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/areas-sociales", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCrearArea_JSONInvalido_400(t *testing.T) {
	r, token := setupRouterSocialAdmin(t)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/areas-sociales", bytes.NewReader([]byte("{esto no es json")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListarAreas_ConToken_200(t *testing.T) {
	r, token := setupRouterSocialAdmin(t)

	body := models.AreaSocial{Nombre: "Salón Comunal", Capacidad: 50}
	buf, _ := json.Marshal(body)
	reqCrear := httptest.NewRequest(http.MethodPost, "/api/v1/areas-sociales", bytes.NewReader(buf))
	reqCrear.Header.Set("Content-Type", "application/json")
	reqCrear.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(httptest.NewRecorder(), reqCrear)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/areas-sociales", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var lista []models.AreaSocial
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &lista))
	assert.Len(t, lista, 1)
	assert.Equal(t, "Salón Comunal", lista[0].Nombre)
}

func TestGetArea_NoExiste_404(t *testing.T) {
	r, token := setupRouterSocial(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/areas-sociales/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetArea_IDInvalido_400(t *testing.T) {
	r, token := setupRouterSocial(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/areas-sociales/abc", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestActualizarArea_Valido_200(t *testing.T) {
	r, token := setupRouterSocialAdmin(t)

	crearBody := models.AreaSocial{Nombre: "Salón Comunal", Capacidad: 50}
	buf, _ := json.Marshal(crearBody)
	reqCrear := httptest.NewRequest(http.MethodPost, "/api/v1/areas-sociales", bytes.NewReader(buf))
	reqCrear.Header.Set("Content-Type", "application/json")
	reqCrear.Header.Set("Authorization", "Bearer "+token)
	recCrear := httptest.NewRecorder()
	r.ServeHTTP(recCrear, reqCrear)
	require.Equal(t, http.StatusCreated, recCrear.Code)

	actualizarBody := models.AreaSocial{Nombre: "Salón Comunal Renovado", Capacidad: 80}
	bufAct, _ := json.Marshal(actualizarBody)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/areas-sociales/1", bytes.NewReader(bufAct))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var respuesta models.AreaSocial
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &respuesta))
	assert.Equal(t, "Salón Comunal Renovado", respuesta.Nombre)
	assert.Equal(t, 80, respuesta.Capacidad)
}

func TestActualizarArea_NombreVacio_400(t *testing.T) {
	r, token := setupRouterSocialAdmin(t)

	buf, _ := json.Marshal(models.AreaSocial{Nombre: "", Capacidad: 10})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/areas-sociales/1", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestBorrarArea_Existente_204(t *testing.T) {
	r, token := setupRouterSocialAdmin(t)

	buf, _ := json.Marshal(models.AreaSocial{Nombre: "Salón Comunal", Capacidad: 50})
	reqCrear := httptest.NewRequest(http.MethodPost, "/api/v1/areas-sociales", bytes.NewReader(buf))
	reqCrear.Header.Set("Content-Type", "application/json")
	reqCrear.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(httptest.NewRecorder(), reqCrear)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/areas-sociales/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestBorrarArea_NoExiste_404(t *testing.T) {
	r, token := setupRouterSocialAdmin(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/areas-sociales/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// ═══════════════════════════════════════════════════════════════════════
// RESERVAS
// ═══════════════════════════════════════════════════════════════════════

func TestCrearReserva_ConToken_201(t *testing.T) {
	// Crear reserva es acción de residente, no solo admin (ver routes_social.go)
	r, token := setupRouterSocial(t)

	inicio := time.Now().Add(24 * time.Hour)
	fin := inicio.Add(2 * time.Hour)
	body := models.ReservaArea{
		AreaID:         1,
		ResidenteID:    1,
		FechaInicio:    inicio,
		FechaFin:       fin,
		Proposito:      "Cumpleaños",
		NumeroPersonas: 15,
	}
	buf, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/reservas-areas", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	var respuesta models.ReservaArea
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &respuesta))
	assert.Equal(t, "Cumpleaños", respuesta.Proposito)
	assert.Equal(t, "pendiente", respuesta.Estado) // el service asigna Estado="pendiente"
	assert.NotZero(t, respuesta.ID)
}

func TestCrearReserva_SinToken_401(t *testing.T) {
	r, _ := setupRouterSocial(t)

	buf, _ := json.Marshal(models.ReservaArea{AreaID: 1, ResidenteID: 1})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/reservas-areas", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestCrearReserva_PropositoVacio_400(t *testing.T) {
	r, token := setupRouterSocial(t)

	inicio := time.Now().Add(24 * time.Hour)
	fin := inicio.Add(2 * time.Hour)
	body := models.ReservaArea{
		AreaID: 1, ResidenteID: 1, Proposito: "",
		NumeroPersonas: 10, FechaInicio: inicio, FechaFin: fin,
	}
	buf, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/reservas-areas", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCrearReserva_NumeroPersonasInvalido_400(t *testing.T) {
	r, token := setupRouterSocial(t)

	inicio := time.Now().Add(24 * time.Hour)
	fin := inicio.Add(2 * time.Hour)
	body := models.ReservaArea{
		AreaID: 1, ResidenteID: 1, Proposito: "Cumpleaños",
		NumeroPersonas: 0, FechaInicio: inicio, FechaFin: fin,
	}
	buf, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/reservas-areas", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCrearReserva_FechasInvalidas_400(t *testing.T) {
	// fecha_fin antes que fecha_inicio → ErrFechasInvalidas
	r, token := setupRouterSocial(t)

	inicio := time.Now().Add(24 * time.Hour)
	fin := inicio.Add(-2 * time.Hour) // inválida a propósito
	body := models.ReservaArea{
		AreaID: 1, ResidenteID: 1, Proposito: "Cumpleaños",
		NumeroPersonas: 10, FechaInicio: inicio, FechaFin: fin,
	}
	buf, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/reservas-areas", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCrearReserva_AreaOResidenteFaltante_404(t *testing.T) {
	// AreaID=0 o ResidenteID=0 → el service devuelve ErrNoEncontrado (404)
	r, token := setupRouterSocial(t)

	inicio := time.Now().Add(24 * time.Hour)
	fin := inicio.Add(2 * time.Hour)
	body := models.ReservaArea{
		AreaID: 0, ResidenteID: 1, Proposito: "Cumpleaños",
		NumeroPersonas: 10, FechaInicio: inicio, FechaFin: fin,
	}
	buf, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/reservas-areas", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestListarReservas_ConToken_200(t *testing.T) {
	r, token := setupRouterSocial(t)

	inicio := time.Now().Add(24 * time.Hour)
	fin := inicio.Add(2 * time.Hour)
	buf, _ := json.Marshal(models.ReservaArea{
		AreaID: 1, ResidenteID: 1, Proposito: "Cumpleaños",
		NumeroPersonas: 10, FechaInicio: inicio, FechaFin: fin,
	})
	reqCrear := httptest.NewRequest(http.MethodPost, "/api/v1/reservas-areas", bytes.NewReader(buf))
	reqCrear.Header.Set("Content-Type", "application/json")
	reqCrear.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(httptest.NewRecorder(), reqCrear)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reservas-areas", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var lista []models.ReservaArea
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &lista))
	assert.Len(t, lista, 1)
}

func TestGetReserva_NoExiste_404(t *testing.T) {
	r, token := setupRouterSocial(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reservas-areas/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestActualizarReserva_TokenNoAdmin_403(t *testing.T) {
	// Actualizar reserva (aprobar/rechazar) es solo admin, ver routes_social.go
	r, token := setupRouterSocial(t)

	buf, _ := json.Marshal(models.ReservaArea{Estado: "aprobada"})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/reservas-areas/1", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestBorrarReserva_ConTokenAdmin_204(t *testing.T) {
	rAdmin, tokenAdmin := setupRouterSocialAdmin(t)

	inicio := time.Now().Add(24 * time.Hour)
	fin := inicio.Add(2 * time.Hour)
	buf, _ := json.Marshal(models.ReservaArea{
		AreaID: 1, ResidenteID: 1, Proposito: "Cumpleaños",
		NumeroPersonas: 10, FechaInicio: inicio, FechaFin: fin,
	})
	reqCrear := httptest.NewRequest(http.MethodPost, "/api/v1/reservas-areas", bytes.NewReader(buf))
	reqCrear.Header.Set("Content-Type", "application/json")
	reqCrear.Header.Set("Authorization", "Bearer "+tokenAdmin)
	rAdmin.ServeHTTP(httptest.NewRecorder(), reqCrear)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/reservas-areas/1", nil)
	req.Header.Set("Authorization", "Bearer "+tokenAdmin)
	rec := httptest.NewRecorder()
	rAdmin.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}

// ═══════════════════════════════════════════════════════════════════════
// NOTIFICACIONES
// ═══════════════════════════════════════════════════════════════════════

func TestCrearNotificacion_ConTokenAdmin_201(t *testing.T) {
	r, token := setupRouterSocialAdmin(t)

	body := models.Notificacion{ResidenteID: 1, Tipo: "aviso", Mensaje: "Corte de agua mañana"}
	buf, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/notificaciones", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	var respuesta models.Notificacion
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &respuesta))
	assert.Equal(t, "Corte de agua mañana", respuesta.Mensaje)
	assert.False(t, respuesta.Leida)
}

func TestCrearNotificacion_TokenNoAdmin_403(t *testing.T) {
	r, token := setupRouterSocial(t)

	buf, _ := json.Marshal(models.Notificacion{ResidenteID: 1, Tipo: "aviso", Mensaje: "Prueba"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/notificaciones", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestCrearNotificacion_MensajeVacio_400(t *testing.T) {
	r, token := setupRouterSocialAdmin(t)

	buf, _ := json.Marshal(models.Notificacion{ResidenteID: 1, Tipo: "aviso", Mensaje: ""})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/notificaciones", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCrearNotificacion_ResidenteIDFaltante_404(t *testing.T) {
	r, token := setupRouterSocialAdmin(t)

	buf, _ := json.Marshal(models.Notificacion{ResidenteID: 0, Tipo: "aviso", Mensaje: "Prueba"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/notificaciones", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestListarNotificaciones_ConToken_200(t *testing.T) {
	r, token := setupRouterSocialAdmin(t)

	buf, _ := json.Marshal(models.Notificacion{ResidenteID: 1, Tipo: "aviso", Mensaje: "Prueba"})
	reqCrear := httptest.NewRequest(http.MethodPost, "/api/v1/notificaciones", bytes.NewReader(buf))
	reqCrear.Header.Set("Content-Type", "application/json")
	reqCrear.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(httptest.NewRecorder(), reqCrear)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/notificaciones", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var lista []models.Notificacion
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &lista))
	assert.Len(t, lista, 1)
}

func TestGetNotificacion_NoExiste_404(t *testing.T) {
	r, token := setupRouterSocial(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/notificaciones/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestMarcarLeida_Existente_200(t *testing.T) {
	r, token := setupRouterSocialAdmin(t)

	buf, _ := json.Marshal(models.Notificacion{ResidenteID: 1, Tipo: "aviso", Mensaje: "Prueba"})
	reqCrear := httptest.NewRequest(http.MethodPost, "/api/v1/notificaciones", bytes.NewReader(buf))
	reqCrear.Header.Set("Content-Type", "application/json")
	reqCrear.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(httptest.NewRecorder(), reqCrear)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/notificaciones/1/leer", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var respuesta models.Notificacion
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &respuesta))
	assert.True(t, respuesta.Leida)
}

func TestMarcarLeida_NoExiste_404(t *testing.T) {
	r, token := setupRouterSocial(t)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/notificaciones/999/leer", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestBorrarNotificacion_ConTokenAdmin_204(t *testing.T) {
	r, token := setupRouterSocialAdmin(t)

	buf, _ := json.Marshal(models.Notificacion{ResidenteID: 1, Tipo: "aviso", Mensaje: "Prueba"})
	reqCrear := httptest.NewRequest(http.MethodPost, "/api/v1/notificaciones", bytes.NewReader(buf))
	reqCrear.Header.Set("Content-Type", "application/json")
	reqCrear.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(httptest.NewRecorder(), reqCrear)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/notificaciones/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestBorrarNotificacion_TokenNoAdmin_403(t *testing.T) {
	r, token := setupRouterSocial(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/notificaciones/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}
