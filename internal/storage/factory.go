package storage

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"condominio-api/internal/models"
)

// Recursos agrupa todo lo que la capa de almacenamiento expone a la aplicación:
// los repositorios ya construidos y una función para cerrar la conexión al apagar.
type Recursos struct {
	Usuarios     UserRepository
	Area         AreaSocialRepository
	Reserva      ReservaAreaRepository // ← corregido: antes decía ReservaRepository
	Notificacion NotificacionRepository
	Cerrar       func() error
}

// Inicializar centraliza todo el plumbing de almacenamiento (patrón Factory):
// abre la conexión SQLite, migra el esquema, siembra datos iniciales y
// construye los repositorios. main.go ya no necesita saber de GORM ni de SQLite.
func Inicializar(rutaDB string) (*Recursos, error) {
	// 1. Abrir conexión y migrar esquema.
	db, err := gorm.Open(sqlite.Open(rutaDB), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("abrir SQLite: %w", err)
	}
	if err := db.AutoMigrate(
		&models.Usuario{},
		&models.AreaSocial{},
		&models.ReservaArea{},
		&models.Notificacion{},
	); err != nil {
		return nil, fmt.Errorf("AutoMigrate: %w", err)
	}

	// 2. Sembrar datos iniciales (solo si están vacíos).
	SembrarSocial(db)

	// 3. Repositorios.
	usuarios := NewUsuarioGORM(db)
	area := NewAreaSQLite(db)
	reserva := NewReservaSQLite(db)
	notificacion := NewNotificacionSQLite(db)

	// 4. Cierre ordenado de la conexión a la base de datos.
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
