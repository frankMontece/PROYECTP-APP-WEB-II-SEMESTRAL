package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"condominio-api/internal/models"
)

// Recursos agrupa los repositorios y la conexión que usará la app.
type Recursos struct {
	DB           *gorm.DB
	Usuarios     UserRepository
	Obligaciones ObligacionRepository
	Multas       MultaRepository
	Cerrar       func() error
}

// Inicializar abre la base de datos según el driver, migra y arma los repos.
func Inicializar(driver, dsn, rutaDB string) (*Recursos, error) {
	gdb, err := abrirGorm(driver, dsn, rutaDB)
	if err != nil {
		return nil, err
	}

	if err := gdb.AutoMigrate(
		&models.Usuario{},
		&models.Obligacion{},
		&models.Multa{},
	); err != nil {
		return nil, fmt.Errorf("automigrate: %w", err)
	}

	usuarioRepo := NewUsuarioGORM(gdb)
	oblRepo := NewObligacionSQLite(gdb)
	multaRepo := NewMultaSQLite(gdb)

	// Sembrar usuarios de ejemplo primero (Obligaciones/Multas dependen de ellos por FK)
	var totalUsuarios int64
	gdb.Model(&models.Usuario{}).Count(&totalUsuarios)
	if totalUsuarios == 0 {
		usuariosEjemplo := []models.Usuario{
			{Email: "residente1@test.com", PasswordHash: "seed", CreadoEn: time.Now()},
			{Email: "residente2@test.com", PasswordHash: "seed", CreadoEn: time.Now()},
			{Email: "residente3@test.com", PasswordHash: "seed", CreadoEn: time.Now()},
		}
		gdb.Create(&usuariosEjemplo)
	}

	// Sembrar datos de ejemplo solo si las tablas están vacías.
	// Útil para tener datos listos en la demo sin crearlos a mano cada vez.
	oblRepo.SembrarVacio()
	multaRepo.SembrarVacio()

	cerrar := func() error {
		sqlDB, err := gdb.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}

	return &Recursos{
		DB:           gdb,
		Usuarios:     usuarioRepo,
		Obligaciones: oblRepo,
		Multas:       multaRepo,
		Cerrar:       cerrar,
	}, nil
}

// abrirGorm decide SQLite o PostgreSQL según DB_DRIVER.
// En postgres reintenta porque el contenedor de la API puede arrancar
// antes de que Postgres esté listo para aceptar conexiones.
func abrirGorm(driver, dsn, rutaDB string) (*gorm.DB, error) {
	switch driver {
	case "postgres":
		var db *gorm.DB
		var err error
		for intento := 1; intento <= 10; intento++ {
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err == nil {
				return db, nil
			}
			log.Printf("PostgreSQL aún no está listo (%d/10): %v", intento, err)
			time.Sleep(2 * time.Second)
		}
		return nil, fmt.Errorf("no fue posible conectar con PostgreSQL: %w", err)

	default: // sqlite
		db, err := gorm.Open(sqlite.Open(rutaDB), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("abrir sqlite: %w", err)
		}
		return db, nil
	}
}
