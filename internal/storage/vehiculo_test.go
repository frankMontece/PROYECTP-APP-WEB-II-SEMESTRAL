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

// nuevaBaseVehiculoTest crea una base SQLite en memoria, migrada y lista
// para usarse en un test aislado — cada llamada devuelve una base nueva,
// sin datos de tests anteriores.
func nuevaBaseVehiculoTest(t *testing.T) *storage.VehiculoSQLite {
	t.Helper()

	gdb, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = gdb.AutoMigrate(&models.Vehiculo{})
	require.NoError(t, err)

	return storage.NewVehiculoSQLite(gdb)
}

func TestSQLiteParqueo_CrearYBuscarVehiculo(t *testing.T) {
	// Preparar
	alm := nuevaBaseVehiculoTest(t)

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
	alm := nuevaBaseVehiculoTest(t)

	// Ejecutar: buscar un ID que no existe
	_, ok := alm.BuscarVehiculoPorID(999)

	// Verificar: debe devolver false (patrón comma-ok)
	require.False(t, ok)
}

func TestSQLiteParqueo_ActualizarVehiculo(t *testing.T) {
	// Preparar: crear un vehículo para luego actualizarlo
	alm := nuevaBaseVehiculoTest(t)

	creado := alm.CrearVehiculo(models.Vehiculo{
		ResidenteID: 1,
		Placa:       "TST-0001",
		Marca:       "Toyota",
		Modelo:      "Corolla",
		Color:       "Blanco",
		Activo:      true,
	})

	nuevosDatos := models.Vehiculo{
		ResidenteID: 1,
		Placa:       "TST-0001",
		Marca:       "Toyota",
		Modelo:      "Corolla Cross", // cambia el modelo
		Color:       "Rojo",          // cambia el color
		Activo:      true,
	}

	// Ejecutar
	actualizado, ok := alm.ActualizarVehiculo(creado.ID, nuevosDatos)

	// Verificar: se aplicó el cambio y el ID no se movió
	require.True(t, ok)
	require.Equal(t, creado.ID, actualizado.ID, "el ID no debe cambiar al actualizar")
	require.Equal(t, "Corolla Cross", actualizado.Modelo)
	require.Equal(t, "Rojo", actualizado.Color)

	// Verificar también contra lo que quedó persistido
	enBase, ok := alm.BuscarVehiculoPorID(creado.ID)
	require.True(t, ok)
	require.Equal(t, "Corolla Cross", enBase.Modelo)
}

func TestSQLiteParqueo_ActualizarVehiculo_NoExiste(t *testing.T) {
	// Preparar: base vacía, ningún vehículo creado
	alm := nuevaBaseVehiculoTest(t)

	// Ejecutar: intentar actualizar un ID que no existe
	_, ok := alm.ActualizarVehiculo(999, models.Vehiculo{
		ResidenteID: 1,
		Placa:       "XXX-0000",
		Marca:       "Marca",
	})

	// Verificar
	require.False(t, ok)
}

func TestSQLiteParqueo_BorrarVehiculo(t *testing.T) {
	// Preparar: crear un vehículo para luego borrarlo
	alm := nuevaBaseVehiculoTest(t)

	creado := alm.CrearVehiculo(models.Vehiculo{
		ResidenteID: 1,
		Placa:       "TST-0001",
		Marca:       "Toyota",
	})

	// Ejecutar
	borrado := alm.BorrarVehiculo(creado.ID)

	// Verificar: se reportó éxito y el registro ya no existe
	require.True(t, borrado)

	_, ok := alm.BuscarVehiculoPorID(creado.ID)
	require.False(t, ok, "el vehículo no debe existir después de borrarlo")
}

func TestSQLiteParqueo_BorrarVehiculo_NoExiste(t *testing.T) {
	// Preparar: base vacía
	alm := nuevaBaseVehiculoTest(t)

	// Ejecutar: borrar un ID que nunca existió
	borrado := alm.BorrarVehiculo(999)

	// Verificar: GORM Delete no reporta error si no hay filas afectadas,
	// así que el repo debe devolver false explícitamente
	require.False(t, borrado)
}

func TestSQLiteParqueo_ListarVehiculos(t *testing.T) {
	// Preparar: crear varios vehículos
	alm := nuevaBaseVehiculoTest(t)

	alm.CrearVehiculo(models.Vehiculo{ResidenteID: 1, Placa: "AAA-0001", Marca: "Toyota"})
	alm.CrearVehiculo(models.Vehiculo{ResidenteID: 2, Placa: "BBB-0002", Marca: "Kia"})
	alm.CrearVehiculo(models.Vehiculo{ResidenteID: 3, Placa: "CCC-0003", Marca: "Chevrolet"})

	// Ejecutar
	lista := alm.ListarVehiculos()

	// Verificar
	require.Len(t, lista, 3)
}

func TestSQLiteParqueo_ListarVehiculos_BaseVacia(t *testing.T) {
	// Preparar: base recién migrada, sin datos
	alm := nuevaBaseVehiculoTest(t)

	// Ejecutar
	lista := alm.ListarVehiculos()

	// Verificar: debe devolver una lista vacía, no nil ni error
	require.Empty(t, lista)
}
