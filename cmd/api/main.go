package main

import (
	"log"
	"net/http"

	"github.com/glebarez/sqlite"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"

	"condominio-api/internal/handlers"
	// línea 13 — el import queda igual (ruta de carpeta)
	middleware "condominio-api/internal/middelware"
	"condominio-api/internal/models"
	"condominio-api/internal/service"
	"condominio-api/internal/storage"
)

func main() {
	// ── 1. Abrir base de datos SQLite con GORM ────────────────────────────
	gdb, err := gorm.Open(sqlite.Open("condominio.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("no se pudo abrir la base de datos: ", err)
	}

	// ── 2. AutoMigrate — crea las tablas si no existen ────────────────────
	if err := gdb.AutoMigrate(
		&models.Usuario{},
		&models.Obligacion{},
		&models.Multa{},
	); err != nil {
		log.Fatal("falló AutoMigrate: ", err)
	}

	// ── 3. Repositorios (capa de storage) ────────────────────────────────
	oblRepo := storage.NewObligacionesGORM(gdb)
	multaRepo := storage.NewMultasGORM(gdb)
	usuarioRepo := storage.NewUsuarioGORM(gdb)

	// ── 4. Services (lógica de negocio) ──────────────────────────────────
	authService := service.NewAuthService(usuarioRepo)
	oblService := service.NewObligacionesService(oblRepo)
	multaService := service.NewMultasService(multaRepo)

	// ── 5. Server (handlers HTTP) ─────────────────────────────────────────
	servidor := handlers.NewServer(authService, oblService, multaService)

	// ── 6. Router ─────────────────────────────────────────────────────────
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.CORS)

	r.Route("/api/v1", func(r chi.Router) {
		// Rutas públicas — no requieren token
		r.Post("/auth/register", servidor.Registrar)
		r.Post("/auth/login", servidor.Login)

		// Rutas protegidas — requieren JWT válido
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authService))

			// Obligaciones
			r.Get("/obligaciones", servidor.ListarObligaciones)
			r.Post("/obligaciones", servidor.CreateObligacion)
			r.Get("/obligaciones/{id}", servidor.GetObligacion)
			r.Put("/obligaciones/{id}", servidor.UpdateObligacion)
			r.Delete("/obligaciones/{id}", servidor.DeleteObligacion)

			// Multas
			r.Get("/multas", servidor.ListarMultas)
			r.Post("/multas", servidor.CreateMulta)
			r.Get("/multas/{id}", servidor.GetMulta)
			r.Put("/multas/{id}", servidor.UpdateMulta)
			r.Delete("/multas/{id}", servidor.DeleteMulta)
		})
	})

	// ── 7. Arrancar servidor ──────────────────────────────────────────────
	log.Println("Servidor escuchando en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
