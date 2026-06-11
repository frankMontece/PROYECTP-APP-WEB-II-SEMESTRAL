package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"condominio-api/internal/handlers"
	"condominio-api/internal/storage"
)

func main() {
	// ── Almacén en memoria del Módulo B ───────────────────────────────────
	storeSocial := storage.NuevaMemoriaSocial()

	// ── Router principal ──────────────────────────────────────────────────
	r := chi.NewRouter()

	// ── Middleware global ─────────────────────────────────────────────────
	// Logger    → imprime método, ruta y duración de cada request
	// Recoverer → si un handler entra en panic devuelve 500 en lugar de
	//             tumbar el servidor completo
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	// ── Rutas versionadas /api/v1 ─────────────────────────────────────────
	r.Route("/api/v1", func(r chi.Router) {

		// Módulo B — Área Social y Notificaciones (Frank Montecé)
		handlers.MontarRutasSocial(r, storeSocial)
	})

	// ── Arrancar el servidor ──────────────────────────────────────────────
	log.Println("=== Módulo B — Área Social y Notificaciones ===")
	log.Println("Servidor escuchando en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
