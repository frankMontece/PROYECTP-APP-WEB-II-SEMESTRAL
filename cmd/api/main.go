package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"condominio-api/internal/config"
	"condominio-api/internal/handlers"
	"condominio-api/internal/httpserver"
	"condominio-api/internal/middleware"
	"condominio-api/internal/service"
	"condominio-api/internal/storage"
)

func main() {
	cfg := config.Cargar()
	if err := run(cfg); err != nil {
		log.Fatal(err)
	}
}

func run(cfg config.Config) error {
	// 1. Inicializar almacenamiento (conexión + migración + siembra + repositorios).
	//    El factory decide sqlite/postgres según cfg.DBDriver, y siembra datos
	//    de ejemplo solo cuando el driver no es postgres.
	recursos, err := storage.Inicializar(cfg.DBDriver, cfg.DBDSN, cfg.RutaDB)
	if err != nil {
		return err
	}
	defer recursos.Cerrar()

	// 2. Servicios de los cuatro módulos.
	authService := service.NewAuthService(
		recursos.Usuarios,
		service.WithSecreto(cfg.JWTSecreto),
		service.WithDuracionToken(cfg.JWTDuracion),
	)

	oblService := service.NewObligacionesService(recursos.Obligaciones)
	multaService := service.NewMultasService(recursos.Multas)

	vehiculoService := service.NewVehiculoService(recursos.Vehiculos)
	visitaService := service.NewVisitaService(recursos.Visitas)
	accesoService := service.NewAccesoService(recursos.Accesos)

	areaService := service.NewAreaService(recursos.Area)
	reservaService := service.NewReservaService(recursos.Reserva)
	notificacionService := service.NewNotificacionService(recursos.Notificacion)

	// 3. Servidor: un solo Services con todos los dominios.
	servidor := handlers.NewServer(handlers.Services{
		Auth: authService,

		Obligaciones: oblService,
		Multas:       multaService,

		Vehiculos: vehiculoService,
		Visitas:   visitaService,
		Accesos:   accesoService,

		Area:         areaService,
		Reserva:      reservaService,
		Notificacion: notificacionService,
	})

	// 4. Router + middleware global.
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.CORS)

	// 5. Rutas /api/v1.
	r.Route("/api/v1", func(r chi.Router) {
		// Rutas públicas.
		r.Post("/auth/register", servidor.Registrar)
		r.Post("/auth/login", servidor.Login)

		// Rutas protegidas — un solo middleware de auth para los cuatro dominios.
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authService))

			handlers.MontarRutasObligaciones(r, servidor)
			handlers.MontarRutasSocial(r, servidor)
			handlers.MontarRutasParqueo(r, servidor)
		})
	})

	// 6. Servidor HTTP con timeouts desde config.
	srv := httpserver.Nuevo(
		r,
		httpserver.ConPuerto(cfg.Puerto),
		httpserver.ConReadTimeout(cfg.ReadTimeout),
		httpserver.ConWriteTimeout(cfg.WriteTimeout),
	)

	// 7. Contexto que se cancela al recibir Ctrl+C o SIGTERM.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 8. Arrancar en goroutine para no bloquear.
	errServidor := make(chan error, 1)
	go func() {
		log.Println("=== Sistema de Gestión de Condominios ===")
		log.Printf("Servidor escuchando en http://localhost%s\n", cfg.Puerto)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errServidor <- err
		}
	}()

	// 9. Esperar señal de apagado o error del servidor.
	select {
	case err := <-errServidor:
		return err
	case <-ctx.Done():
		log.Println("Señal de apagado recibida, cerrando ordenadamente...")
	}

	// 10. Graceful shutdown — hasta 10s para terminar requests en curso.
	ctxApagado, cancelar := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelar()

	if err := srv.Shutdown(ctxApagado); err != nil {
		return err
	}

	log.Println("Servidor detenido limpiamente.")
	return nil
}
