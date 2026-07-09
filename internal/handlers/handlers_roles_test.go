package handlers_test

// setupRouterAdmin es igual a setupRouter, pero el usuario de prueba se
// inserta directamente en el fake con Rol: "admin" (sin pasar por
// AuthService.Registrar, que siempre asigna "residente"). Se usa en los
// tests de escritura (POST/PUT/DELETE), que ahora exigen rol admin.
//
// setupRouter (el original) se mantiene tal cual y sigue sirviendo para
// los tests de solo lectura y para el caso "sin token" / "rol insuficiente".

import (
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"condominio-api/internal/handlers"
	"condominio-api/internal/middleware"
	"condominio-api/internal/models"
	"condominio-api/internal/service"
)

func setupRouterAdmin(t *testing.T) (chi.Router, string) {
	t.Helper()

	// Generamos el hash real ANTES de construir el fake, para poder hacer
	// Login() con la misma contraseña en texto plano usada aquí.
	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	require.NoError(t, err)

	repoUsuarios := &fakeUserRepo{
		usuarios: []models.Usuario{
			{ID: 1, Email: "admin@test.com", PasswordHash: string(hash), Rol: service.RolAdmin},
		},
	}
	authService := service.NewAuthService(repoUsuarios)

	token, err := authService.Login("admin@test.com", "password123")
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
