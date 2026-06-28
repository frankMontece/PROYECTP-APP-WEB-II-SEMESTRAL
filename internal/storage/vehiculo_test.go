package storage_test

// Patrón: test del repositorio real con GORM contra sqlite :memory:
// sqlite.Open(":memory:") crea una base en RAM que se descarta al terminar el test.
// No hay mocks aquí: probamos la persistencia real.
// Si GORM no asigna ID, si la query de búsqueda está mal, este test lo detecta.

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"condominio-api/internal/models"
	"condominio-api/internal/storage"
)

func TestSQLiteParqueo_CrearYBuscarVehiculo(t *testing.T) {
	// Preparar: SQLite en memoria, desechable
	// Cada vez que corre el test arranca desde cero (sin datos previos)
	gdb, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = gdb.AutoMigrate(&models.Vehiculo{})
	require.NoError(t, err)

	alm := storage.NewSQLiteParqueo(gdb)

	entrada := models.Vehiculo{
		ResidenteID: 1,
		Placa:       "TST-0001",
		Marca:       "Toyota",
		Modelo:      "Corolla",
		Color:       "Blanco",
		Activo:      true,
	}

	// Ejecutar: crear
	creado := alm.CrearVehiculo(entrada)

	// Verificar: GORM asignó el ID automáticamente
	require.NotZero(t, creado.ID, "GORM debe asignar un ID distinto de 0")

	// Ejecutar: buscar por el ID recién creado
	encontrado, ok := alm.BuscarVehiculoPorID(creado.ID)

	// Verificar: el registro refleja lo que se guardó
	require.True(t, ok, "BuscarVehiculoPorID debe encontrar el registro recién creado")
	require.Equal(t, "TST-0001", encontrado.Placa)
	require.Equal(t, "Toyota", encontrado.Marca)
	require.Equal(t, uint(1), encontrado.ResidenteID)
}

func TestSQLiteParqueo_BuscarVehiculo_NoExiste(t *testing.T) {
	// Preparar
	gdb, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	gdb.AutoMigrate(&models.Vehiculo{})

	alm := storage.NewSQLiteParqueo(gdb)

	// Ejecutar: buscar un ID que no existe
	_, ok := alm.BuscarVehiculoPorID(999)

	// Verificar: debe devolver false (patrón comma-ok)
	require.False(t, ok)
}
