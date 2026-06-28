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
