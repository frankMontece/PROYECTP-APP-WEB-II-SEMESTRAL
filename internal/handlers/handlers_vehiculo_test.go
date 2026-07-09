package handlers_test

// Patrón: handler con httptest + fake en memoria

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"condominio-api/internal/handlers"
	"condominio-api/internal/middleware"
	"condominio-api/internal/models"
	"condominio-api/internal/service"
)

// ── Fake del repositorio de usuarios ─────────────────────────────────────────

type fakeUserRepo struct {
	usuarios []models.Usuario
}

func (f *fakeUserRepo) CrearUsuario(u models.Usuario) (models.Usuario, error) {
	u.ID = len(f.usuarios) + 1
	f.usuarios = append(f.usuarios, u)
	return u, nil
}

func (f *fakeUserRepo) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	for _, u := range f.usuarios {
		if u.Email == email {
			return u, true
		}
	}
	return models.Usuario{}, false
}

// ── Fake del almacén de parqueo (implementa storage.AlmacenParqueo) ─────────

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

// ── Helper: construye router + devuelve token de prueba ──────────────────────

func setupRouter(t *testing.T) (chi.Router, string) {
	t.Helper()

	// Auth
	repoUsuarios := &fakeUserRepo{}
	authService := service.NewAuthService(repoUsuarios)

	_, err := authService.Registrar("test@test.com", "password123")
	require.NoError(t, err)

	token, err := authService.Login("test@test.com", "password123")
	require.NoError(t, err)

	// Parqueo (un solo almacén implementa las 3 interfaces)
	almacen := &fakeAlmacenParqueo{}
	vehiculoSvc := service.NewVehiculoService(almacen)
	visitaSvc := service.NewVisitaService(almacen)
	accesoSvc := service.NewAccesoService(almacen)

	srv := handlers.NewServer(handlers.Services{
		Auth:      authService,
		Vehiculos: vehiculoSvc,
		Visitas:   visitaSvc,
		Accesos:   accesoSvc,
	})

	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(authService))
		handlers.MontarRutasParqueo(r, srv)
	})

	return r, token
}

// ── Test 1: POST /vehiculos con token → 201 ──────────────────────────────────

func TestCrearVehiculo_ConToken_201(t *testing.T) {
	// Preparar
	r, token := setupRouterAdmin(t)

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

// ── Test 2: POST /vehiculos SIN token → 401 ──────────────────────────────────

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

// ── Test 3: GET /vehiculos/{id} inexistente → 404 ───────────────────────────

func TestObtenerVehiculo_NoExiste_404(t *testing.T) {
	r, token := setupRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/vehiculos/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// ── Test 4: GET /vehiculos con resultados → 200 ──────────────────────────────

func TestListarVehiculos_ConToken_200(t *testing.T) {
	// Preparar
	r, token := setupRouterAdmin(t)

	// Sembrar un vehículo antes de listar, para no depender de datos vacíos
	body := models.Vehiculo{ResidenteID: 1, Placa: "PBG-2241", Marca: "Toyota"}
	buf, _ := json.Marshal(body)
	reqCrear := httptest.NewRequest(http.MethodPost, "/vehiculos", bytes.NewReader(buf))
	reqCrear.Header.Set("Content-Type", "application/json")
	reqCrear.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(httptest.NewRecorder(), reqCrear)

	req := httptest.NewRequest(http.MethodGet, "/vehiculos", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	// Ejecutar
	r.ServeHTTP(rec, req)

	// Verificar
	require.Equal(t, http.StatusOK, rec.Code)

	var lista []models.Vehiculo
	err := json.Unmarshal(rec.Body.Bytes(), &lista)
	require.NoError(t, err)
	assert.Len(t, lista, 1)
	assert.Equal(t, "PBG-2241", lista[0].Placa)
}

// ── Test 5: PUT /vehiculos/{id} válido → 200 ─────────────────────────────────

func TestActualizarVehiculo_Valido_200(t *testing.T) {
	// Preparar: crear un vehículo para luego actualizarlo
	r, token := setupRouterAdmin(t)

	crearBody := models.Vehiculo{ResidenteID: 1, Placa: "PBG-2241", Marca: "Toyota", Modelo: "Corolla"}
	buf, _ := json.Marshal(crearBody)
	reqCrear := httptest.NewRequest(http.MethodPost, "/vehiculos", bytes.NewReader(buf))
	reqCrear.Header.Set("Content-Type", "application/json")
	reqCrear.Header.Set("Authorization", "Bearer "+token)
	recCrear := httptest.NewRecorder()
	r.ServeHTTP(recCrear, reqCrear)

	var creado models.Vehiculo
	require.NoError(t, json.Unmarshal(recCrear.Body.Bytes(), &creado))

	// Ejecutar: actualizar el modelo y color
	actualizarBody := models.Vehiculo{
		ResidenteID: 1,
		Placa:       "PBG-2241",
		Marca:       "Toyota",
		Modelo:      "Corolla Cross",
		Color:       "Rojo",
	}
	bufAct, _ := json.Marshal(actualizarBody)
	req := httptest.NewRequest(http.MethodPut, "/vehiculos/1", bytes.NewReader(bufAct))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	// Verificar
	require.Equal(t, http.StatusOK, rec.Code)

	var respuesta models.Vehiculo
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &respuesta))
	assert.Equal(t, "Corolla Cross", respuesta.Modelo)
	assert.Equal(t, "Rojo", respuesta.Color)
}

// ── Test 6: PUT /vehiculos/{id} con body inválido → 400 ──────────────────────

func TestActualizarVehiculo_BodyInvalido_400(t *testing.T) {
	// Preparar
	r, token := setupRouterAdmin(t)

	req := httptest.NewRequest(http.MethodPut, "/vehiculos/1", bytes.NewReader([]byte("{esto no es json")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	// Ejecutar
	r.ServeHTTP(rec, req)

	// Verificar
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ── Test 7: DELETE /vehiculos/{id} no existe → 404 ───────────────────────────

func TestBorrarVehiculo_NoExiste_404(t *testing.T) {
	// Preparar: base vacía, ningún vehículo creado
	r, token := setupRouterAdmin(t)

	req := httptest.NewRequest(http.MethodDelete, "/vehiculos/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	// Ejecutar
	r.ServeHTTP(rec, req)

	// Verificar: el service devuelve ErrVehiculoNoEncontrado → statusDeError → 404
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// ── Test 8: DELETE /vehiculos/{id} existente → 204 ───────────────────────────

func TestBorrarVehiculo_Existente_204(t *testing.T) {
	// Preparar: crear un vehículo para luego borrarlo
	r, token := setupRouterAdmin(t)

	crearBody := models.Vehiculo{ResidenteID: 1, Placa: "PBG-2241", Marca: "Toyota"}
	buf, _ := json.Marshal(crearBody)
	reqCrear := httptest.NewRequest(http.MethodPost, "/vehiculos", bytes.NewReader(buf))
	reqCrear.Header.Set("Content-Type", "application/json")
	reqCrear.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(httptest.NewRecorder(), reqCrear)

	// Ejecutar
	req := httptest.NewRequest(http.MethodDelete, "/vehiculos/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	// Verificar
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

// ── Test 9: POST /vehiculos con JSON inválido → 400 ──────────────────────────

func TestCrearVehiculo_JSONInvalido_400(t *testing.T) {
	// Preparar
	r, token := setupRouterAdmin(t)

	req := httptest.NewRequest(http.MethodPost, "/vehiculos", bytes.NewReader([]byte("{esto no es json")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	// Ejecutar
	r.ServeHTTP(rec, req)

	// Verificar
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ── Test 10: POST /vehiculos con campo requerido vacío → 400 ─────────────────

func TestCrearVehiculo_PlacaVacia_400(t *testing.T) {
	// Preparar
	r, token := setupRouterAdmin(t)

	body := models.Vehiculo{ResidenteID: 1, Placa: "", Marca: "Toyota"}
	buf, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/vehiculos", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	// Ejecutar
	r.ServeHTTP(rec, req)

	// Verificar: el service rechaza la placa vacía antes de llegar al repo
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
