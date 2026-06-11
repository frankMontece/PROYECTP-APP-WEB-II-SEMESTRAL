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
	// ── 1. Almacenes en memoria (uno por módulo) ──────────────────────────
	//storeParqueo := storage.NuevaMemoriaParqueo()
	storeAlicuotas := storage.NuevaMemoriaObligaciones()

	// ── 2. Router principal ───────────────────────────────────────────────
	r := chi.NewRouter()

	// ── 3. Middleware global ──────────────────────────────────────────────
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	// ── 4. Rutas versionadas /api/v1 ─────────────────────────────────────
	r.Route("/api/v1", func(r chi.Router) {

		// Módulo C — Parqueo y Vehículos (Winter Povea)
		//handlers.MontarRutasParqueo(r, storeParqueo)

		// Módulo A — Obligaciones (Héctor Fernández)
		handlers.MontarRutasAlicuotas(r, storeAlicuotas)

		// Módulo B — Área Social (Frank Montecé) → pendiente de merge
		// handlers.MontarRutasSocial(r, storeSocial)
	})

	// ── 5. Arrancar el servidor ───────────────────────────────────────────
	log.Println("Servidor escuchando en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
