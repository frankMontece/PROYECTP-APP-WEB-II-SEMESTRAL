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
	if err := db.AutoMigrate(&models.Usuario{}, &models.Obligacion{}); err != nil {
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

// TestObligacionesGORM_BuscarPorID_NoExiste verifica que BuscarObligacionPorID
// devuelve ok=false cuando el ID no existe en la BD.
func TestObligacionesGORM_BuscarPorID_NoExiste(t *testing.T) {
	db := abrirDBPrueba(t)
	repo := storage.NewObligacionSQLite(db)

	_, ok := repo.BuscarObligacionPorID(999)
	if ok {
		t.Fatal("esperaba ok=false para un ID que no existe, obtuvo true")
	}
}

// TestObligacionesGORM_Actualizar_Existente verifica que ActualizarObligacion
// persiste los cambios cuando el registro existe.
func TestObligacionesGORM_Actualizar_Existente(t *testing.T) {
	db := abrirDBPrueba(t)
	repo := storage.NewObligacionSQLite(db)

	creada := repo.CrearObligacion(models.Obligacion{
		ResidenteID: 1,
		Tipo:        "mensual",
		Monto:       100.00,
		Periodo:     "2026-06",
		Estado:      "pendiente",
	})

	datosNuevos := models.Obligacion{
		ResidenteID: 1,
		Tipo:        "mensual",
		Monto:       150.00,
		Periodo:     "2026-06",
		Estado:      "pagada",
	}

	actualizada, ok := repo.ActualizarObligacion(creada.ID, datosNuevos)
	if !ok {
		t.Fatalf("ActualizarObligacion devolvió false para ID %d que sí existe", creada.ID)
	}
	if actualizada.Monto != 150.00 {
		t.Fatalf("monto esperado 150.00, obtenido %.2f", actualizada.Monto)
	}
	if actualizada.Estado != "pagada" {
		t.Fatalf("estado esperado 'pagada', obtenido '%s'", actualizada.Estado)
	}

	// Confirmar que el cambio se persistió en BD, no solo en memoria
	confirmada, ok := repo.BuscarObligacionPorID(creada.ID)
	if !ok {
		t.Fatal("no se encontró la obligación tras actualizar")
	}
	if confirmada.Estado != "pagada" {
		t.Fatalf("el cambio no se persistió: estado sigue siendo '%s'", confirmada.Estado)
	}
}

// TestObligacionesGORM_Actualizar_NoExiste verifica que ActualizarObligacion
// devuelve ok=false cuando el ID no existe.
func TestObligacionesGORM_Actualizar_NoExiste(t *testing.T) {
	db := abrirDBPrueba(t)
	repo := storage.NewObligacionSQLite(db)

	_, ok := repo.ActualizarObligacion(999, models.Obligacion{
		ResidenteID: 1,
		Tipo:        "mensual",
		Monto:       50.00,
		Periodo:     "2026-06",
	})

	if ok {
		t.Fatal("esperaba ok=false al actualizar un ID que no existe, obtuvo true")
	}
}

// TestObligacionesGORM_Borrar_Existente verifica que BorrarObligacion
// elimina el registro y devuelve true.
func TestObligacionesGORM_Borrar_Existente(t *testing.T) {
	db := abrirDBPrueba(t)
	repo := storage.NewObligacionSQLite(db)

	creada := repo.CrearObligacion(models.Obligacion{
		ResidenteID: 1,
		Tipo:        "mensual",
		Monto:       80.00,
		Periodo:     "2026-06",
	})

	ok := repo.BorrarObligacion(creada.ID)
	if !ok {
		t.Fatalf("BorrarObligacion devolvió false para ID %d que sí existe", creada.ID)
	}

	// Confirmar que ya no se encuentra
	_, existe := repo.BuscarObligacionPorID(creada.ID)
	if existe {
		t.Fatal("la obligación seguía existiendo después de borrarla")
	}
}

// TestObligacionesGORM_Borrar_NoExiste verifica que BorrarObligacion
// devuelve false cuando el ID no existe.
func TestObligacionesGORM_Borrar_NoExiste(t *testing.T) {
	db := abrirDBPrueba(t)
	repo := storage.NewObligacionSQLite(db)

	ok := repo.BorrarObligacion(999)
	if ok {
		t.Fatal("esperaba false al borrar un ID que no existe, obtuvo true")
	}
}

// TestObligacionesGORM_SembrarVacio_InsertaCuandoEstaVacia verifica que
// SembrarVacio inserta datos de ejemplo cuando la tabla está vacía.
func TestObligacionesGORM_SembrarVacio_InsertaCuandoEstaVacia(t *testing.T) {
	db := abrirDBPrueba(t)
	repo := storage.NewObligacionSQLite(db)

	storage.SembrarObligacion(db)

	lista := repo.ListarObligaciones()
	if len(lista) != 3 {
		t.Fatalf("esperaba 3 obligaciones sembradas, hay %d", len(lista))
	}
}

// TestObligacionesGORM_SembrarVacio_NoDuplicaSiYaHayDatos verifica que
// SembrarVacio no inserta datos si la tabla ya tiene registros.
func TestObligacionesGORM_SembrarVacio_NoDuplicaSiYaHayDatos(t *testing.T) {
	db := abrirDBPrueba(t)
	repo := storage.NewObligacionSQLite(db)

	// Insertamos manualmente un registro antes de sembrar
	repo.CrearObligacion(models.Obligacion{
		ResidenteID: 1,
		Tipo:        "mensual",
		Monto:       50.00,
		Periodo:     "2026-06",
	})

	storage.SembrarObligacion(db) // no debería hacer nada, ya hay 1 registro

	lista := repo.ListarObligaciones()
	if len(lista) != 1 {
		t.Fatalf("esperaba que SembrarObligacion no insertara nada (ya había 1 registro), hay %d", len(lista))
	}
}
