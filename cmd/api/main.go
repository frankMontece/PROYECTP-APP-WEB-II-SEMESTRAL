package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"

	"condominio-api/internal/config"
	"condominio-api/internal/handlers"
	"condominio-api/internal/httpserver"
	middleware "condominio-api/internal/middelware"
	"condominio-api/internal/models"
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
	// 1. Abrir base de datos SQLite con GORM
	gdb, err := gorm.Open(sqlite.Open(cfg.RutaDB), &gorm.Config{})
	if err != nil {
		return err
	}

	// 2. AutoMigrate — crea las tablas si no existen
	if err := gdb.AutoMigrate(
		&models.Usuario{},
		&models.Obligacion{},
		&models.Multa{},
	); err != nil {
		return err
	}

	// 3. Repositorios (capa de storage)
	oblRepo := storage.NewObligacionesGORM(gdb)
	multaRepo := storage.NewMultasGORM(gdb)
	usuarioRepo := storage.NewUsuarioGORM(gdb)

	// 4. Services (lógica de negocio)
	// El secreto JWT ahora viene de config, ya no es una variable global
	authService := service.NewAuthService(usuarioRepo, cfg.JWTSecreto, cfg.JWTDuracion)
	oblService := service.NewObligacionesService(oblRepo)
	multaService := service.NewMultasService(multaRepo)

	// 5. Server (handlers HTTP)
	servidor := handlers.NewServer(authService, oblService, multaService)

	// 6. Router + middleware global
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.CORS)

	r.Route("/api/v1", func(r chi.Router) {
		// Rutas públicas
		r.Post("/auth/register", servidor.Registrar)
		r.Post("/auth/login", servidor.Login)

		// Rutas protegidas con JWT
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authService))

			r.Get("/obligaciones", servidor.ListarObligaciones)
			r.Post("/obligaciones", servidor.CreateObligacion)
			r.Get("/obligaciones/{id}", servidor.GetObligacion)
			r.Put("/obligaciones/{id}", servidor.UpdateObligacion)
			r.Delete("/obligaciones/{id}", servidor.DeleteObligacion)

			r.Get("/multas", servidor.ListarMultas)
			r.Post("/multas", servidor.CreateMulta)
			r.Get("/multas/{id}", servidor.GetMulta)
			r.Put("/multas/{id}", servidor.UpdateMulta)
			r.Delete("/multas/{id}", servidor.DeleteMulta)
		})
	})

	// 7. Servidor HTTP con timeouts desde config
	srv := httpserver.Nuevo(
		r,
		httpserver.ConPuerto(cfg.Puerto),
		httpserver.ConReadTimeout(cfg.ReadTimeout),
		httpserver.ConWriteTimeout(cfg.WriteTimeout),
	)

	// 8. Contexto que se cancela al recibir Ctrl+C o SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 9. Arrancar en goroutine para no bloquear
	errServidor := make(chan error, 1)
	go func() {
		log.Printf("Servidor escuchando en http://localhost%s", cfg.Puerto)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errServidor <- err
		}
	}()

	// 10. Esperar señal de apagado o error del servidor
	select {
	case err := <-errServidor:
		return err
	case <-ctx.Done():
		log.Println("Señal de apagado recibida, cerrando ordenadamente...")
	}

	// 11. Graceful shutdown — hasta 10s para terminar requests en curso
	ctxApagado, cancelar := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelar()
	if err := srv.Shutdown(ctxApagado); err != nil {
		return err
	}
	log.Println("Servidor detenido limpiamente.")
	return nil
}
