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

// Recursos agrupa todo lo que la capa de almacenamiento expone a la aplicación:
// los repositorios ya construidos y una función para cerrar la conexión al apagar.
type Recursos struct {
	DB *gorm.DB

	Usuarios UserRepository

	Obligaciones ObligacionRepository
	Multas       MultaRepository

	Vehiculos VehiculosRepository
	Visitas   VisitasRepository
	Accesos   AccesosRepository

	Area         AreaSocialRepository
	Reserva      ReservaAreaRepository
	Notificacion NotificacionRepository

	Cerrar func() error
}

// Inicializar centraliza todo el plumbing de almacenamiento (patrón Factory):
// abre la conexión (SQLite o PostgreSQL según driver), migra el esquema,
// siembra datos iniciales y construye los repositorios. main.go ya no necesita
// saber de GORM ni del driver concreto.
//
// driver: "sqlite" (default) o "postgres". dsn: solo se usa si driver == "postgres".
// rutaDB: solo se usa si driver == "sqlite" (o está vacío/default).
func Inicializar(driver, dsn, rutaDB string) (*Recursos, error) {
	gdb, err := abrirGorm(driver, dsn, rutaDB)
	if err != nil {
		return nil, err
	}

	// 2. Migrar esquema (idéntico, sin importar el driver).
	if err := gdb.AutoMigrate(
		&models.Usuario{},
		&models.Obligacion{},
		&models.Multa{},
		&models.Vehiculo{},
		&models.VisitaVehiculo{},
		&models.AccesoVehiculo{},
		&models.AreaSocial{},
		&models.ReservaArea{},
		&models.Notificacion{},
	); err != nil {
		return nil, fmt.Errorf("AutoMigrate: %w", err)
	}

	// 3. Repositorios.
	usuarios := NewUsuarioGORM(gdb)
	obligaciones := NewObligacionSQLite(gdb)
	multas := NewMultaSQLite(gdb)
	vehiculos := NewVehiculoSQLite(gdb)
	visitas := NewVisitaSQLite(gdb)
	accesos := NewAccesoSQLite(gdb)
	area := NewAreaSQLite(gdb)
	reserva := NewReservaSQLite(gdb)
	notificacion := NewNotificacionSQLite(gdb)

	// 4. Sembrar datos iniciales — siempre, incluyendo postgres, para tener
	//    datos de ejemplo listos en la demo. Cada función es idempotente:
	//    solo inserta si la tabla correspondiente está vacía.
	SembrarUsuarios(gdb)
	SembrarSocial(gdb)
	SembrarObligacion(gdb)
	SembrarParqueo(gdb)

	// 5. Cierre ordenado de la conexión a la base de datos.
	cerrar := func() error {
		sqlDB, err := gdb.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}

	return &Recursos{
		DB: gdb,

		Usuarios: usuarios,

		Obligaciones: obligaciones,
		Multas:       multas,

		Vehiculos: vehiculos,
		Visitas:   visitas,
		Accesos:   accesos,

		Area:         area,
		Reserva:      reserva,
		Notificacion: notificacion,

		Cerrar: cerrar,
	}, nil
}

// abrirGorm decide automáticamente si usar SQLite o PostgreSQL,
// con reintentos en postgres porque el contenedor de la API puede
// arrancar antes de que el contenedor de Postgres esté listo.
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
