package storage

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"condominio-api/internal/models"
)

// Recursos agrupa todo lo que la capa de almacenamiento expone a la aplicación:
// los repositorios ya construidos y una función para cerrar la conexión al apagar.
type Recursos struct {
	Usuarios     UserRepository
	Area         AreaSocialRepository
	Reserva      ReservaAreaRepository
	Notificacion NotificacionRepository
	Cerrar       func() error
}

// Inicializar centraliza todo el plumbing de almacenamiento (patrón Factory):
// abre la conexión (SQLite o PostgreSQL según driver), migra el esquema,
// siembra datos iniciales y construye los repositorios. main.go ya no necesita
// saber de GORM ni del driver concreto.
//
// driver: "sqlite" (default) o "postgres". dsn: solo se usa si driver == "postgres".
// rutaDB: solo se usa si driver == "sqlite" (o está vacío/default).
func Inicializar(driver, rutaDB, dsn string) (*Recursos, error) {
	// 1. Abrir conexión según el driver. ÚNICO CAMBIO DE LÓGICA DEL HITO.
	var (
		db  *gorm.DB
		err error
	)
	switch driver {
	case "postgres":
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("abrir PostgreSQL: %w", err)
		}
	default: // "sqlite" o vacío
		db, err = gorm.Open(sqlite.Open(rutaDB), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("abrir SQLite: %w", err)
		}
	}

	// 2. Migrar esquema (idéntico, sin importar el driver).
	if err := db.AutoMigrate(
		&models.Usuario{},
		&models.AreaSocial{},
		&models.ReservaArea{},
		&models.Notificacion{},
	); err != nil {
		return nil, fmt.Errorf("AutoMigrate: %w", err)
	}

	// 3. Sembrar datos iniciales (solo si están vacíos).
	SembrarSocial(db)

	// 4. Repositorios.
	usuarios := NewUsuarioGORM(db)
	area := NewAreaSQLite(db)
	reserva := NewReservaSQLite(db)
	notificacion := NewNotificacionSQLite(db)

	// 5. Cierre ordenado de la conexión a la base de datos.
	cerrar := func() error {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}

	return &Recursos{
		Usuarios:     usuarios,
		Area:         area,
		Reserva:      reserva,
		Notificacion: notificacion,
		Cerrar:       cerrar,
	}, nil
}
