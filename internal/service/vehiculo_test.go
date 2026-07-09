package service_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"condominio-api/internal/models"
	"condominio-api/internal/service"
)

// mockVehiculoRepo implementa storage.VehiculosRepository con testify/mock.
// Cada método llama a m.Called(...) para registrar que fue invocado
// y devolver lo que configuramos con On(...).Return(...).

type mockVehiculoRepo struct {
	mock.Mock
}

func (m *mockVehiculoRepo) ListarVehiculos() []models.Vehiculo {
	args := m.Called()
	return args.Get(0).([]models.Vehiculo)
}

func (m *mockVehiculoRepo) BuscarVehiculoPorID(id uint) (models.Vehiculo, bool) {
	args := m.Called(id)
	return args.Get(0).(models.Vehiculo), args.Bool(1)
}

func (m *mockVehiculoRepo) CrearVehiculo(v models.Vehiculo) models.Vehiculo {
	args := m.Called(v)
	return args.Get(0).(models.Vehiculo)
}

func (m *mockVehiculoRepo) ActualizarVehiculo(id uint, datos models.Vehiculo) (models.Vehiculo, bool) {
	args := m.Called(id, datos)
	return args.Get(0).(models.Vehiculo), args.Bool(1)
}

func (m *mockVehiculoRepo) BorrarVehiculo(id uint) bool {
	args := m.Called(id)
	return args.Bool(0)
}

func TestVehiculoService_Crear_Validaciones(t *testing.T) {
	casos := []struct {
		nombre    string
		entrada   models.Vehiculo
		esperaErr error
	}{
		{
			nombre:    "placa vacía debe rechazarse",
			entrada:   models.Vehiculo{ResidenteID: 1, Placa: "", Marca: "Toyota"},
			esperaErr: service.ErrPlacaRequerida,
		},
		{
			nombre:    "residente_id cero debe rechazarse",
			entrada:   models.Vehiculo{ResidenteID: 0, Placa: "PBG-2241", Marca: "Toyota"},
			esperaErr: service.ErrResidenteRequerido,
		},
		{
			nombre:    "marca vacía debe rechazarse",
			entrada:   models.Vehiculo{ResidenteID: 1, Placa: "PBG-2241", Marca: ""},
			esperaErr: service.ErrMarcaRequerida,
		},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			// Preparar: repo mock sin ninguna configuración de On/Return
			// porque esperamos que nunca sea llamado
			repo := new(mockVehiculoRepo)
			svc := service.NewVehiculoService(repo)

			// Ejecutar
			_, err := svc.Crear(c.entrada)

			// Verificar
			require.ErrorIs(t, err, c.esperaErr)
			// La regla de negocio falló ANTES de llegar al repositorio
			repo.AssertNotCalled(t, "CrearVehiculo")
		})
	}
}

func TestVehiculoService_Crear_Valido(t *testing.T) {
	// Preparar
	repo := new(mockVehiculoRepo)
	esperado := models.Vehiculo{ID: 1, ResidenteID: 1, Placa: "PBG-2241", Marca: "Toyota", Activo: true}
	// Cuando alguien llame CrearVehiculo con cualquier argumento, devolver `esperado`
	repo.On("CrearVehiculo", mock.Anything).Return(esperado)

	svc := service.NewVehiculoService(repo)

	entrada := models.Vehiculo{
		ResidenteID: 1,
		Placa:       "PBG-2241",
		Marca:       "Toyota",
		Modelo:      "Corolla",
		Color:       "Blanco",
	}

	// Ejecutar
	creado, err := svc.Crear(entrada)

	// Verificar
	require.NoError(t, err)
	require.Equal(t, uint(1), creado.ID)
	require.True(t, creado.Activo) // el service asigna Activo=true
	repo.AssertCalled(t, "CrearVehiculo", mock.Anything)
}

func TestVehiculoService_Obtener_Existente(t *testing.T) {
	// Preparar
	repo := new(mockVehiculoRepo)
	esperado := models.Vehiculo{ID: 1, ResidenteID: 1, Placa: "PBG-2241", Marca: "Toyota", Activo: true}
	repo.On("BuscarVehiculoPorID", uint(1)).Return(esperado, true)

	svc := service.NewVehiculoService(repo)

	// Ejecutar
	v, err := svc.Obtener(1)

	// Verificar
	require.NoError(t, err)
	require.Equal(t, esperado, v)
	repo.AssertCalled(t, "BuscarVehiculoPorID", uint(1))
}

func TestVehiculoService_Obtener_NoExistente(t *testing.T) {
	// Preparar
	repo := new(mockVehiculoRepo)
	repo.On("BuscarVehiculoPorID", uint(999)).Return(models.Vehiculo{}, false)

	svc := service.NewVehiculoService(repo)

	// Ejecutar
	_, err := svc.Obtener(999)

	// Verificar
	require.ErrorIs(t, err, service.ErrVehiculoNoEncontrado)
}

// ---------- Actualizar ----------

func TestVehiculoService_Actualizar_Valido(t *testing.T) {
	// Preparar
	repo := new(mockVehiculoRepo)
	datos := models.Vehiculo{ResidenteID: 1, Placa: "PBG-2241", Marca: "Toyota", Modelo: "Corolla"}
	actualizado := models.Vehiculo{ID: 1, ResidenteID: 1, Placa: "PBG-2241", Marca: "Toyota", Modelo: "Corolla"}
	repo.On("ActualizarVehiculo", uint(1), datos).Return(actualizado, true)

	svc := service.NewVehiculoService(repo)

	// Ejecutar
	v, err := svc.Actualizar(1, datos)

	// Verificar
	require.NoError(t, err)
	require.Equal(t, actualizado, v)
	repo.AssertCalled(t, "ActualizarVehiculo", uint(1), datos)
}

func TestVehiculoService_Actualizar_NoExistente(t *testing.T) {
	// Preparar
	repo := new(mockVehiculoRepo)
	datos := models.Vehiculo{ResidenteID: 1, Placa: "PBG-2241", Marca: "Toyota"}
	repo.On("ActualizarVehiculo", uint(999), datos).Return(models.Vehiculo{}, false)

	svc := service.NewVehiculoService(repo)

	// Ejecutar
	_, err := svc.Actualizar(999, datos)

	// Verificar
	require.ErrorIs(t, err, service.ErrVehiculoNoEncontrado)
}

func TestVehiculoService_Actualizar_DatosInvalidos(t *testing.T) {
	// Preparar: placa vacía — la validación debe fallar ANTES de llegar al repo
	repo := new(mockVehiculoRepo)
	svc := service.NewVehiculoService(repo)

	datos := models.Vehiculo{ResidenteID: 1, Placa: "", Marca: "Toyota"}

	// Ejecutar
	_, err := svc.Actualizar(1, datos)

	// Verificar
	require.ErrorIs(t, err, service.ErrPlacaRequerida)
	repo.AssertNotCalled(t, "ActualizarVehiculo")
}

// ---------- Borrar ----------

func TestVehiculoService_Borrar_Valido(t *testing.T) {
	// Preparar
	repo := new(mockVehiculoRepo)
	repo.On("BorrarVehiculo", uint(1)).Return(true)

	svc := service.NewVehiculoService(repo)

	// Ejecutar
	err := svc.Borrar(1)

	// Verificar
	require.NoError(t, err)
	repo.AssertCalled(t, "BorrarVehiculo", uint(1))
}

func TestVehiculoService_Borrar_NoExistente(t *testing.T) {
	// Preparar
	repo := new(mockVehiculoRepo)
	repo.On("BorrarVehiculo", uint(999)).Return(false)

	svc := service.NewVehiculoService(repo)

	// Ejecutar
	err := svc.Borrar(999)

	// Verificar
	require.ErrorIs(t, err, service.ErrVehiculoNoEncontrado)
}

// ---------- Listar ----------

func TestVehiculoService_Listar_DevuelveTodos(t *testing.T) {
	// Preparar
	repo := new(mockVehiculoRepo)
	esperado := []models.Vehiculo{
		{ID: 1, Placa: "PBG-2241"},
		{ID: 2, Placa: "ACG-0873"},
	}
	repo.On("ListarVehiculos").Return(esperado)

	svc := service.NewVehiculoService(repo)

	// Ejecutar
	lista := svc.Listar()

	// Verificar
	require.Len(t, lista, 2)
	require.Equal(t, esperado, lista)
}

// mock.Anything se usa arriba en el test original de Crear; lo dejamos
// disponible aquí también por si se necesita en futuros casos.
var _ = mock.Anything
