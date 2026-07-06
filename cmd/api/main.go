package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"condominio-api/internal/config"
	"condominio-api/internal/handlers"
	"condominio-api/internal/middleware"
	"condominio-api/internal/service"
	"condominio-api/internal/storage"
)

func main() {
	cfg := config.Cargar()

	recursos, err := storage.Inicializar(
		cfg.DBDriver,
		cfg.DBDSN,
		cfg.RutaDB,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer recursos.Cerrar()

	authService := service.NewAuthService(
		recursos.Usuarios,
		cfg.JWTSecreto,
		cfg.JWTDuracion,
	)

	vehiculoService := service.NewVehiculoService(recursos.Parqueo)
	visitaService := service.NewVisitaService(recursos.Parqueo)
	accesoService := service.NewAccesoService(recursos.Parqueo)

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
