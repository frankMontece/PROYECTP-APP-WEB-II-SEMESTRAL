package main

import (
	"context"
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
	// 0. Cargar configuracion (.env + defaults).
	cfg := config.Cargar()

	// 1. Inicializar almacenamiento (conexión + migración + siembra + repositorios).
	//    main.go ya no importa GORM ni ningún driver concreto: eso vive en storage.Inicializar.
	recursos, err := storage.Inicializar(cfg.DBDriver, cfg.RutaDB, cfg.DBDSN)
	if err != nil {
		log.Fatal("Error al inicializar el almacenamiento:", err)
	}

	// 2. Servicios. AuthService recibe secreto/duracion desde Config via Options.
	authService := service.NewAuthService(
		recursos.Usuarios,
		service.WithSecreto(cfg.JWTSecreto),
		service.WithDuracionToken(cfg.JWTDuracion),
	)
	areaService := service.NewAreaService(recursos.Area)
	reservaService := service.NewReservaService(recursos.Reserva)
	notificacionService := service.NewNotificacionService(recursos.Notificacion)

	// 3. Servidor (struct de dependencias, Bloque 4).
	servidor := handlers.NewServer(handlers.Deps{
		Auth:         authService,
		Area:         areaService,
		Reserva:      reservaService,
		Notificacion: notificacionService,
	})

	// 4. Router.
	r := chi.NewRouter()

	// 5. Middleware global.
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.CORS)

	// 6. Rutas /api/v1.
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", servidor.Registrar)
		r.Post("/auth/login", servidor.Login)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authService))
			handlers.MontarRutasSocial(r, servidor)
		})
	})

	// 7. http.Server con timeouts configurables.
	srv := httpserver.Nuevo(r,
		httpserver.ConPuerto(cfg.Puerto),
		httpserver.ConReadTimeout(cfg.ReadTimeout),
		httpserver.ConWriteTimeout(cfg.WriteTimeout),
	)

	// 8. Contexto que se cancela al recibir Ctrl+C / SIGTERM.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 9. Arrancar el servidor en una goroutine.
	go func() {
		log.Println("=== Módulo B — Área Social y Notificaciones (Hito 3) ===")
		log.Printf("🚀 Servidor escuchando en http://localhost%s\n", cfg.Puerto)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error en el servidor: %v", err)
		}
	}()

	// 10. Esperar señal de apagado.
	<-ctx.Done()
	log.Println("Señal de apagado recibida, cerrando ordenadamente...")

	// 11. Dar hasta 10s a las peticiones en curso.
	ctxApagado, cancelar := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelar()

	if err := srv.Shutdown(ctxApagado); err != nil {
		log.Fatalf("Error al apagar el servidor: %v", err)
	}

	// 12. Cerrar la conexión a la base de datos (antes de terminar el proceso).
	if err := recursos.Cerrar(); err != nil {
		log.Printf("Error al cerrar la base de datos: %v", err)
	}

	log.Println("Servidor apagado correctamente.")
}
