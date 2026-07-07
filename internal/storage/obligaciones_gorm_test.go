package storage_test

//El test que prueba que esa conexión con la base de datos real funciona, crea un registro, lo busca, lo encuentra.
import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"condominio-api/internal/models"
	"condominio-api/internal/storage"
)

// abrirDBPrueba abre una BD SQLite en RAM y corre AutoMigrate.
// Cada test llama a esta función y obtiene una BD limpia e independiente.
func abrirDBPrueba(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("no se pudo abrir SQLite en memoria: %v", err)
	}
	if err := db.AutoMigrate(&models.Obligacion{}); err != nil {
		t.Fatalf("AutoMigrate falló: %v", err)
	}
	return db
}

// TestObligacionesGORM_CrearYBuscar verifica que:
// 1. GORM asigna un ID al crear la obligación (persiste en BD).
// 2. BuscarObligacionPorID la encuentra con los datos correctos.
// 3. ListarObligaciones la refleja en el slice resultado.
func TestObligacionesGORM_CrearYBuscar(t *testing.T) {
	db := abrirDBPrueba(t)
	repo := storage.NewObligacionSQLite(db)

	// Crear
	nueva := models.Obligacion{
		ResidenteID: 1,
		Tipo:        "mensual",
		Monto:       200.00,
		Periodo:     "2026-06",
		Estado:      "pendiente",
	}
	creada := repo.CrearObligacion(nueva)

	// GORM debe haber asignado un ID
	if creada.ID == 0 {
		t.Fatal("GORM debe asignar un ID al crear la obligación, obtuvo 0")
	}

	// Buscar por ID — debe encontrarse
	encontrada, ok := repo.BuscarObligacionPorID(creada.ID)
	if !ok {
		t.Fatalf("BuscarObligacionPorID devolvió false para ID %d que sí existe", creada.ID)
	}
	if encontrada.Monto != 200.00 {
		t.Fatalf("monto esperado 200.00, obtenido %.2f", encontrada.Monto)
	}
	if encontrada.Periodo != "2026-06" {
		t.Fatalf("periodo esperado '2026-06', obtenido '%s'", encontrada.Periodo)
	}

	// Listar — debe aparecer exactamente 1
	lista := repo.ListarObligaciones()
	if len(lista) != 1 {
		t.Fatalf("esperado 1 obligación en la lista, hay %d", len(lista))
	}
}
