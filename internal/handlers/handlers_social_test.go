package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"condominio-api/internal/middleware"
	"condominio-api/internal/models"
	"condominio-api/internal/service"
	"condominio-api/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── Fake en memoria (doble simple, sin mock) ──────────────────────────────────

type fakeAreaRepo struct {
	areas []models.AreaSocial
}

func (f *fakeAreaRepo) ListarAreas() []models.AreaSocial { return f.areas }
func (f *fakeAreaRepo) BuscarAreaPorID(id uint) (models.AreaSocial, bool) {
	for _, a := range f.areas {
		if a.ID == id {
			return a, true
		}
	}
	return models.AreaSocial{}, false
}
func (f *fakeAreaRepo) CrearArea(a models.AreaSocial) models.AreaSocial {
	a.ID = uint(len(f.areas) + 1)
	f.areas = append(f.areas, a)
	return a
}
func (f *fakeAreaRepo) ActualizarArea(id uint, datos models.AreaSocial) (models.AreaSocial, bool) {
	for i, a := range f.areas {
		if a.ID == id {
			datos.ID = id
			f.areas[i] = datos
			return datos, true
		}
	}
	return models.AreaSocial{}, false
}
func (f *fakeAreaRepo) BorrarArea(id uint) bool {
	for i, a := range f.areas {
		if a.ID == id {
			f.areas = append(f.areas[:i], f.areas[i+1:]...)
			return true
		}
	}
	return false
}

// Verifica en tiempo de compilación que fakeAreaRepo cumple la interfaz
var _ storage.AreaSocialRepository = (*fakeAreaRepo)(nil)

// ── Helper: construye router con middleware Auth ───────────────────────────────

func nuevoRouterDeTest(t *testing.T) (chi.Router, *service.AuthService) {
	t.Helper()

	fakeUserRepo := &fakeUserRepo{}
	authSvc := service.NewAuthService(fakeUserRepo)

	areaRepo := &fakeAreaRepo{
		areas: []models.AreaSocial{
			{
				ID:        1,
				Nombre:    "Salón Principal",
				Capacidad: 100,
				Activo:    true,
			},
		},
	}

	areaSvc := service.NewAreaService(areaRepo)

	srv := NewServer(Services{
		Auth: authSvc,
		Area: areaSvc,
	})

	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authSvc))
			MontarRutasSocial(r, srv)
		})
	})

	return r, authSvc
}

// ── Fake UserRepo para AuthService ────────────────────────────────────────────

type fakeUserRepo struct {
	usuarios []models.Usuario
}

var _ storage.UserRepository = (*fakeUserRepo)(nil)

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

// ── Tests ─────────────────────────────────────────────────────────────────────

// Test A: ruta protegida SIN token → debe responder 401
func TestListarAreas_SinToken_Responde401(t *testing.T) {
	router, _ := nuevoRouterDeTest(t)

	// Preparar — request sin Authorization
	req := httptest.NewRequest(http.MethodGet, "/api/v1/areas-sociales", nil)
	rec := httptest.NewRecorder()

	// Ejecutar
	router.ServeHTTP(rec, req)

	// Verificar
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

// Test B: ruta protegida CON token válido → debe responder 200 y devolver áreas
func TestListarAreas_ConToken_Responde200(t *testing.T) {
	router, authSvc := nuevoRouterDeTest(t)

	// Registrar un usuario y obtener token
	_, err := authSvc.Registrar("test@test.com", "password123")
	require.NoError(t, err)
	token, err := authSvc.Login("test@test.com", "password123")
	require.NoError(t, err)

	// Preparar — request con token
	req := httptest.NewRequest(http.MethodGet, "/api/v1/areas-sociales", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	// Ejecutar
	router.ServeHTTP(rec, req)

	// Verificar
	assert.Equal(t, http.StatusOK, rec.Code)

	var areas []models.AreaSocial
	err = json.NewDecoder(rec.Body).Decode(&areas)
	require.NoError(t, err)
	assert.Len(t, areas, 1)
	assert.Equal(t, "Salón Principal", areas[0].Nombre)
}
