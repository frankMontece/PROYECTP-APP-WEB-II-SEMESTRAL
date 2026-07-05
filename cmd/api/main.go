package main

import (
	"log"
	"net/http"

	"github.com/glebarez/sqlite"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"

	"condominio-api/internal/config"
	"condominio-api/internal/handlers"
	"condominio-api/internal/middleware"
	"condominio-api/internal/models"
	"condominio-api/internal/service"
	"condominio-api/internal/storage"
)

func main() {
	// 1. Cargar configuración (.env + variables de entorno)
	cfg := config.Cargar()

	// 2. Conectar a SQLite y migrar esquema.
	db, err := gorm.Open(sqlite.Open(cfg.RutaDB), &gorm.Config{})
	if err != nil {
		log.Fatal("Error al conectar a la base de datos:", err)
	}

	if err := db.AutoMigrate(
		&models.Usuario{},
		&models.Vehiculo{},
		&models.VisitaVehiculo{},
		&models.AccesoVehiculo{},
	); err != nil {
		log.Fatal("Error al migrar la base de datos:", err)
	}

	// 3. Repositorios e insertar datos iniciales.
	userRepo := storage.NewUsuarioGORM(db)
	almacen := storage.NewSQLiteParqueo(db)
	almacen.SembrarVacio()

	// 4. Servicios.
	authService := service.NewAuthService(
		userRepo,
		cfg.JWTSecreto,
		cfg.JWTDuracion,
	)

	vehiculoService := service.NewVehiculoService(almacen)
	visitaService := service.NewVisitaService(almacen)
	accesoService := service.NewAccesoService(almacen)

	// 5. Servidor (handlers HTTP).
	services := handlers.Services{
		Auth:      authService,
		Vehiculos: vehiculoService,
		Visitas:   visitaService,
		Accesos:   accesoService,
	}

	servidor := handlers.NewServer(services)

	// 6. Router.
	r := chi.NewRouter()

	// 7. Middleware global.
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.CORS)

	// 8. Rutas /api/v1.
	r.Route("/api/v1", func(r chi.Router) {

		// Rutas públicas
		servidor.MontarRutasAuth(r)

		// Rutas protegidas
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authService))

			servidor.MontarRutasParqueo(r)
		})
	})

	// 9. Configuración del servidor HTTP
	server := &http.Server{
		Addr:         cfg.Puerto,
		Handler:      r,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	log.Printf("Servidor escuchando en http://localhost%s", cfg.Puerto)

	log.Fatal(server.ListenAndServe())
}
