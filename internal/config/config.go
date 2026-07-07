package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config agrupa toda la configuracion del servidor en un solo lugar.
//
// Antes: el secreto JWT vivia en una var global de service/auth.go, el puerto
// y la ruta de la DB eran literales en main.go. Ahora hay UNA sola fuente de verdad.
type Config struct {
	Puerto       string
	RutaDB       string
	DBDriver     string // "sqlite" (default, local) o "postgres" (Docker)
	DBDSN        string // DSN completo de Postgres; vacío si DBDriver == "sqlite"
	JWTSecreto   []byte
	JWTDuracion  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func Cargar() Config {
	_ = godotenv.Load()

	return Config{
		Puerto:       conTexto("PUERTO", ":8080"),
		RutaDB:       conTexto("RUTA_DB", "condominio.db"),
		DBDriver:     conTexto("DB_DRIVER", "sqlite"),
		DBDSN:        conTexto("DB_DSN", ""),
		JWTSecreto:   []byte(conTexto("JWT_SECRETO", "condominio-secreto-solo-dev")),
		JWTDuracion:  conDuracion("JWT_DURACION", 24*time.Hour),
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
