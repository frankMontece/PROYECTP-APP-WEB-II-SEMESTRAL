package storage

import (
	"testing"

	"condominio-api/internal/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestAreaSQLite_CrearYBuscar(t *testing.T) {
	// Preparar — base de datos desechable en memoria
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("no se pudo abrir SQLite :memory: → %v", err)
	}
	db.AutoMigrate(&models.AreaSocial{})

	repo := NewAreaSQLite(db)

	entrada := models.AreaSocial{
		Nombre:    "Sala de Reuniones",
		Capacidad: 20,
		Activo:    true,
	}

	// Ejecutar
	creada := repo.CrearArea(entrada)

	// Verificar — el ID fue asignado por GORM
	if creada.ID == 0 {
		t.Fatal("GORM no asignó ID al área creada")
	}

	encontrada, ok := repo.BuscarAreaPorID(creada.ID)
	if !ok {
		t.Fatalf("BuscarAreaPorID(%d) devolvió false", creada.ID)
	}
	if encontrada.Nombre != entrada.Nombre {
		t.Errorf("Nombre esperado %q, obtenido %q", entrada.Nombre, encontrada.Nombre)
	}
	if encontrada.Capacidad != entrada.Capacidad {
		t.Errorf("Capacidad esperada %d, obtenida %d", entrada.Capacidad, encontrada.Capacidad)
	}
}

func TestAreaSQLite_ListarAreas(t *testing.T) {
	// Preparar
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{}) // glebarez: sin CGO
	db.AutoMigrate(&models.AreaSocial{})
	repo := NewAreaSQLite(db)

	repo.CrearArea(models.AreaSocial{Nombre: "Área A", Capacidad: 10, Activo: true})
	repo.CrearArea(models.AreaSocial{Nombre: "Área B", Capacidad: 30, Activo: true})

	// Ejecutar
	areas := repo.ListarAreas()

	// Verificar
	if len(areas) != 2 {
		t.Errorf("esperaba 2 áreas, obtuvo %d", len(areas))
	}
}
