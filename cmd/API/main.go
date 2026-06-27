package main

import (
	"log"
	"net/http"

	"github.com/glebarez/sqlite"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"

	"condominio-api/internal/handlers"
	"condominio-api/internal/middleware"
	"condominio-api/internal/models"
	"condominio-api/internal/service"
	"condominio-api/internal/storage"
)

func main() {
	// ── 1. Conectar a SQLite ──────────────────────────────────────────────
	gdb, err := gorm.Open(sqlite.Open("condominio.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("no se pudo abrir la base de datos: ", err)
	}
	log.Println("Conectado a SQLite")

	// ── 2. AutoMigrate ─────────────────────────────────────────────────────
	if err := gdb.AutoMigrate(
		&models.Usuario{},
		&models.AreaSocial{},
		&models.ReservaArea{},
		&models.Notificacion{},
	); err != nil {
		log.Fatal("falló AutoMigrate: ", err)
	}
	log.Println("Migraciones aplicadas")

	// ── 3. Repositorios (Storage) ──────────────────────────────────────────
	socialRepo := storage.NewSocialGORM(gdb)
	usuarioRepo := storage.NewUsuarioGORM(gdb)

	// ── 4. Servicios (Service) ─────────────────────────────────────────────
	authService := service.NewAuthService(usuarioRepo)
	socialService := service.NewSocialService(socialRepo, socialRepo, socialRepo)

	// ── 5. Server (Handlers) ──────────────────────────────────────────────
	servidor := handlers.NewServer(socialService, authService)

	// ── 6. Router ──────────────────────────────────────────────────────────
	r := chi.NewRouter()

	// Middleware global
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.CORS)

	// ── 7. Rutas ──────────────────────────────────────────────────────────
	r.Route("/api/v1", func(r chi.Router) {
		// Rutas públicas (sin autenticación)
		r.Post("/auth/register", servidor.Registrar)
		r.Post("/auth/login", servidor.Login)

		// Rutas protegidas (requieren JWT)
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authService))

			// Módulo B - Área Social y Notificaciones
			handlers.MontarRutasSocial(r, servidor)
		})
	})

	// ── 8. Arrancar el servidor ──────────────────────────────────────────
	log.Println("=== Módulo B — Área Social y Notificaciones (Hito 3) ===")
	log.Println("Servidor escuchando en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
