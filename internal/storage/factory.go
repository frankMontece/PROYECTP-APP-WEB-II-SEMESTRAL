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

// Recursos agrupa todos los repositorios que utilizará la aplicación.
type Recursos struct {
	DB       *gorm.DB
	Usuarios UserRepository
	Parqueo  *SQLiteParqueo
	Cerrar   func() error
}

// Inicializar abre la base de datos indicada por la configuración,
// ejecuta las migraciones, siembra datos y devuelve los repositorios.
func Inicializar(driver, dsn, rutaDB string) (*Recursos, error) {

	gdb, err := abrirGorm(driver, dsn, rutaDB)
	if err != nil {
		return nil, err
	}

	// Migraciones
	err = gdb.AutoMigrate(
		&models.Usuario{},
		&models.Vehiculo{},
		&models.VisitaVehiculo{},
		&models.AccesoVehiculo{},
	)

	if err != nil {
		return nil, fmt.Errorf("automigrate: %w", err)
	}

	// Repositorios
	userRepo := NewUsuarioGORM(gdb)

	almacen := NewSQLiteParqueo(gdb)
	almacen.SembrarVacio()

	// Función para cerrar la conexión
	cerrar := func() error {
		sqlDB, err := gdb.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}

	return &Recursos{
		DB:       gdb,
		Usuarios: userRepo,
		Parqueo:  almacen,
		Cerrar:   cerrar,
	}, nil
}

// abrirGorm decide automáticamente si usar SQLite o PostgreSQL.
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

			log.Printf(
				"PostgreSQL aún no está listo (%d/10): %v",
				intento,
				err,
			)

			time.Sleep(2 * time.Second)
		}

		return nil, fmt.Errorf("no fue posible conectar con PostgreSQL: %w", err)

	default:

		db, err := gorm.Open(
			sqlite.Open(rutaDB),
			&gorm.Config{},
		)

		if err != nil {
			return nil, fmt.Errorf("abrir sqlite: %w", err)
		}

		return db, nil
	}
}
