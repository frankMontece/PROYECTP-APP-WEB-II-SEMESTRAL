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
	// 1. Conectar a SQLite y migrar esquema.
	db, err := gorm.Open(sqlite.Open("condominio.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Error al conectar a la base de datos:", err)
	}
	if err := db.AutoMigrate(
		&models.Usuario{},
		&models.AreaSocial{},
		&models.ReservaArea{},
		&models.Notificacion{},
	); err != nil {
		log.Fatal("Error al migrar la base de datos:", err)
	}

	// 2. Sembrar datos iniciales (solo si están vacíos).
	storage.SembrarSocial(db) // ← Función separada

	// 3. Repositorios (cada uno con su propia implementación SQLite).
	userRepo := storage.NewUsuarioGORM(db)
	areaRepo := storage.NewAreaSQLite(db)                 // ← Repositorio de áreas
	reservaRepo := storage.NewReservaSQLite(db)           // ← Repositorio de reservas
	notificacionRepo := storage.NewNotificacionSQLite(db) // ← Repositorio de notificaciones

	// 4. Servicios (lógica de negocio + validaciones) - SEGMENTADOS.
	authService := service.NewAuthService(userRepo)
	areaService := service.NewAreaService(areaRepo)
	reservaService := service.NewReservaService(reservaRepo)
	notificacionService := service.NewNotificacionService(notificacionRepo)

	// 5. Servidor (handlers HTTP).
	servidor := handlers.NewServer(
		authService,
		areaService,
		reservaService,
		notificacionService,
	)

	// 6. Router.
	r := chi.NewRouter()

	// 7. Middleware global.
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.CORS)

	// 8. Rutas /api/v1.
	r.Route("/api/v1", func(r chi.Router) {
		// Rutas públicas
		r.Post("/auth/register", servidor.Registrar)
		r.Post("/auth/login", servidor.Login)

		// Rutas protegidas con JWT
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authService))

			// Módulo B - Área Social y Notificaciones
			handlers.MontarRutasSocial(r, servidor)
		})
	})

	log.Println("=== Módulo B — Área Social y Notificaciones (Hito 3) ===")
	log.Println("🚀 Servidor escuchando en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
