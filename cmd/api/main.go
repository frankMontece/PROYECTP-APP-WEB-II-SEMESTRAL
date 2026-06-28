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
		&models.Vehiculo{},
		&models.VisitaVehiculo{},
		&models.AccesoVehiculo{},
	); err != nil {
		log.Fatal("Error al migrar la base de datos:", err)
	}

	// 2. Repositorios e insertar datos iniciales.
	userRepo := storage.NewUsuarioGORM(db)
	almacen := storage.NewSQLiteParqueo(db)
	almacen.SembrarVacio()

	// 3. Servicios (lógica de negocio + validaciones).
	authService := service.NewAuthService(userRepo)
	vehiculoService := service.NewVehiculoService(almacen)
	visitaService := service.NewVisitaService(almacen)
	accesoService := service.NewAccesoService(almacen)

	// 4. Servidor (handlers HTTP).
	servidor := handlers.NewServer(authService, vehiculoService, visitaService, accesoService)

	// 5. Router.
	r := chi.NewRouter()

	// 6. Middleware global.
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.CORS)

	// 7. Rutas /api/v1.
	r.Route("/api/v1", func(r chi.Router) {
		servidor.MontarRutasAuth(r)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authService))

			servidor.MontarRutasParqueo(r)
		})
	})

	log.Println("Servidor escuchando en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
