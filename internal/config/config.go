// Package config carga la configuración de la aplicación desde variables de
// entorno (con soporte para un archivo .env opcional) y expone valores por
// defecto razonables para desarrollo.
package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config agrupa toda la configuración del servidor en un solo lugar.
// Antes estos valores estaban dispersos y hardcodeados: el secreto JWT
// vivía en una var global de service/auth.go, el puerto y la ruta de la
// DB eran literales en main.go. Ahora hay UNA sola fuente de verdad.
type Config struct {
	Puerto       string        // puerto HTTP, ej ":8080"
	RutaDB       string        // archivo SQLite, ej "condominio.db"
	JWTSecreto   []byte        // clave para firmar/verificar JWT
	JWTDuracion  time.Duration // validez del token
	ReadTimeout  time.Duration // timeout de lectura del servidor HTTP
	WriteTimeout time.Duration // timeout de escritura del servidor HTTP
}

// Cargar lee la configuración desde .env o variables de entorno del sistema.
// Si no hay .env, no es un error: en producción las variables vienen del entorno real.
func Cargar() Config {
	_ = godotenv.Load()

	return Config{
		Puerto:       conTexto("PUERTO", ":8080"),
		RutaDB:       conTexto("RUTA_DB", "condominio.db"),
		JWTSecreto:   []byte(conTexto("JWT_SECRETO", "condominio-dev-secret-2026")),
		JWTDuracion:  conDuracion("JWT_DURACION", 2*time.Hour),
		ReadTimeout:  conDuracion("HTTP_READ_TIMEOUT", 10*time.Second),
		WriteTimeout: conDuracion("HTTP_WRITE_TIMEOUT", 10*time.Second),
	}
}

func conTexto(clave, porDefecto string) string {
	if v := os.Getenv(clave); v != "" {
		return v
	}
	return porDefecto
}

func conDuracion(clave string, porDefecto time.Duration) time.Duration {
	v := os.Getenv(clave)
	if v == "" {
		return porDefecto
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return porDefecto
	}
	return d
}
