package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Puerto       string
	DBDriver     string
	DBDSN        string
	RutaDB       string
	JWTSecreto   []byte
	JWTDuracion  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func Cargar() Config {
	_ = godotenv.Load()

	return Config{
		Puerto:       conTexto("PUERTO", ":8080"),
		DBDriver:     conTexto("DB_DRIVER", "sqlite3"),
		DBDSN:        conTexto("DB_DSN", ""),
		RutaDB:       conTexto("RUTA_DB", "condominio.db"),
		JWTSecreto:   []byte(conTexto("JWT_SECRETO", "condominio-dev-secret-2026")),
		JWTDuracion:  conDuracion("JWT_DURACION", 24*time.Hour),
		ReadTimeout:  conDuracion("HTTP_READ_TIMEOUT", 10*time.Second),
		WriteTimeout: conDuracion("HTTP_WRITE_TIMEOUT", 10*time.Second),
	}
}

// conTexto devuelve la variable de entorno o el valor por defecto si esta vacia.
func conTexto(clave, porDefecto string) string {
	if v := os.Getenv(clave); v != "" {
		return v
	}
	return porDefecto
}

// conDuracion parsea una duracion (ej "24h", "30m") o usa el valor por defecto.
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
