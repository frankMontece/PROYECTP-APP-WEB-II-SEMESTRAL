package service

import (
	"testing"

	"condominio-api/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockObligacionRepo implementa storage.ObligacionRepository usando testify/mock.
type mockObligacionRepo struct {
	mock.Mock
}

func (m *mockObligacionRepo) ListarObligaciones() []models.Obligacion {
	args := m.Called()
	return args.Get(0).([]models.Obligacion)
}

func (m *mockObligacionRepo) BuscarObligacionPorID(id uint) (models.Obligacion, bool) {
	args := m.Called(id)
	return args.Get(0).(models.Obligacion), args.Bool(1)
}

func (m *mockObligacionRepo) CrearObligacion(o models.Obligacion) models.Obligacion {
	args := m.Called(o)
	return args.Get(0).(models.Obligacion)
}

func (m *mockObligacionRepo) ActualizarObligacion(id uint, datos models.Obligacion) (models.Obligacion, bool) {
	args := m.Called(id, datos)
	return args.Get(0).(models.Obligacion), args.Bool(1)
}

func (m *mockObligacionRepo) BorrarObligacion(id uint) bool {
	args := m.Called(id)
	return args.Bool(0)
}

// ---------- CrearObligacion ----------

func TestCrearObligacion_RechazaResidenteIDInvalido(t *testing.T) {
	mockRepo := new(mockObligacionRepo)
	svc := NewObligacionesService(mockRepo)

	_, err := svc.CrearObligacion(models.Obligacion{
		ResidenteID: 0, // inválido
		Tipo:        "mensual",
		Monto:       50.00,
		Periodo:     "2026-06",
	})

	assert.Equal(t, ErrResidenteIDInvalido, err)
	mockRepo.AssertNotCalled(t, "CrearObligacion")
}

func TestCrearObligacion_RechazaMontoInvalido(t *testing.T) {
	mockRepo := new(mockObligacionRepo)
	svc := NewObligacionesService(mockRepo)

	_, err := svc.CrearObligacion(models.Obligacion{
		ResidenteID: 1,
		Tipo:        "mensual",
		Monto:       -50.00, // inválido
		Periodo:     "2026-06",
	})

	assert.Equal(t, ErrMontoInvalido, err)
	mockRepo.AssertNotCalled(t, "CrearObligacion")
}

func TestCrearObligacion_RechazaPeriodoVacio(t *testing.T) {
	mockRepo := new(mockObligacionRepo)
	svc := NewObligacionesService(mockRepo)

	_, err := svc.CrearObligacion(models.Obligacion{
		ResidenteID: 1,
		Tipo:        "mensual",
		Monto:       50.00,
		Periodo:     "   ", // inválido
	})

	assert.Equal(t, ErrPeriodoVacio, err)
	mockRepo.AssertNotCalled(t, "CrearObligacion")
}

func TestCrearObligacion_RechazaTipoInvalido(t *testing.T) {
	mockRepo := new(mockObligacionRepo)
	svc := NewObligacionesService(mockRepo)

	_, err := svc.CrearObligacion(models.Obligacion{
		ResidenteID: 1,
		Tipo:        "otro", // inválido
		Monto:       50.00,
		Periodo:     "2026-06",
	})

	assert.Equal(t, ErrTipoObligacion, err)
	mockRepo.AssertNotCalled(t, "CrearObligacion")
}

func TestCrearObligacion_Exitosa(t *testing.T) {
	mockRepo := new(mockObligacionRepo)
	svc := NewObligacionesService(mockRepo)

	entrada := models.Obligacion{
		ResidenteID: 1,
		Tipo:        "mensual",
		Monto:       85.50,
		Periodo:     "2026-06",
	}
	esperado := entrada
	esperado.Estado = "pendiente"
	esperado.ID = 1

	// El service setea Estado="pendiente" antes de llamar al repo,
	// así que el mock debe esperar el objeto ya con ese cambio.
	entradaConEstado := entrada
	entradaConEstado.Estado = "pendiente"

	mockRepo.On("CrearObligacion", entradaConEstado).Return(esperado)

	resultado, err := svc.CrearObligacion(entrada)

	assert.NoError(t, err)
	assert.Equal(t, "pendiente", resultado.Estado)
	assert.Equal(t, uint(1), resultado.ID)
	mockRepo.AssertExpectations(t)
}

// ---------- ObtenerObligacion ----------

func TestObtenerObligacion_Encontrada(t *testing.T) {
	mockRepo := new(mockObligacionRepo)
	svc := NewObligacionesService(mockRepo)

	esperado := models.Obligacion{ID: 5, ResidenteID: 1, Monto: 85.50}
	mockRepo.On("BuscarObligacionPorID", uint(5)).Return(esperado, true)

	resultado, err := svc.ObtenerObligacion(5)

	assert.NoError(t, err)
	assert.Equal(t, esperado, resultado)
	mockRepo.AssertExpectations(t)
}

func TestObtenerObligacion_NoEncontrada(t *testing.T) {
	mockRepo := new(mockObligacionRepo)
	svc := NewObligacionesService(mockRepo)

	mockRepo.On("BuscarObligacionPorID", uint(99)).Return(models.Obligacion{}, false)

	_, err := svc.ObtenerObligacion(99)

	assert.Equal(t, ErrNoEncontrado, err)
	mockRepo.AssertExpectations(t)
}

// ---------- ActualizarObligacion ----------

func TestActualizarObligacion_RechazaMontoInvalido(t *testing.T) {
	mockRepo := new(mockObligacionRepo)
	svc := NewObligacionesService(mockRepo)

	_, err := svc.ActualizarObligacion(1, models.Obligacion{
		ResidenteID: 1,
		Tipo:        "mensual",
		Monto:       0, // inválido
		Periodo:     "2026-06",
	})

	assert.Equal(t, ErrMontoInvalido, err)
	mockRepo.AssertNotCalled(t, "ActualizarObligacion")
}

func TestActualizarObligacion_NoEncontrada(t *testing.T) {
	mockRepo := new(mockObligacionRepo)
	svc := NewObligacionesService(mockRepo)

	datos := models.Obligacion{ResidenteID: 1, Tipo: "mensual", Monto: 50, Periodo: "2026-06"}
	mockRepo.On("ActualizarObligacion", uint(1), datos).Return(models.Obligacion{}, false)

	_, err := svc.ActualizarObligacion(1, datos)

	assert.Equal(t, ErrNoEncontrado, err)
	mockRepo.AssertExpectations(t)
}

func TestActualizarObligacion_Exitosa(t *testing.T) {
	mockRepo := new(mockObligacionRepo)
	svc := NewObligacionesService(mockRepo)

	datos := models.Obligacion{ResidenteID: 1, Tipo: "mensual", Monto: 50, Periodo: "2026-06"}
	esperado := datos
	esperado.ID = 1

	mockRepo.On("ActualizarObligacion", uint(1), datos).Return(esperado, true)

	resultado, err := svc.ActualizarObligacion(1, datos)

	assert.NoError(t, err)
	assert.Equal(t, esperado, resultado)
	mockRepo.AssertExpectations(t)
}

// ---------- BorrarObligacion ----------

func TestBorrarObligacion_Exitosa(t *testing.T) {
	mockRepo := new(mockObligacionRepo)
	svc := NewObligacionesService(mockRepo)

	mockRepo.On("BorrarObligacion", uint(1)).Return(true)

	err := svc.BorrarObligacion(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestBorrarObligacion_NoEncontrada(t *testing.T) {
	mockRepo := new(mockObligacionRepo)
	svc := NewObligacionesService(mockRepo)

	mockRepo.On("BorrarObligacion", uint(99)).Return(false)

	err := svc.BorrarObligacion(99)

	assert.Equal(t, ErrNoEncontrado, err)
	mockRepo.AssertExpectations(t)
}

// ---------- ListarObligaciones ----------

func TestListarObligaciones(t *testing.T) {
	mockRepo := new(mockObligacionRepo)
	svc := NewObligacionesService(mockRepo)

	esperado := []models.Obligacion{{ID: 1}, {ID: 2}}
	mockRepo.On("ListarObligaciones").Return(esperado)

	resultado := svc.ListarObligaciones()

	assert.Equal(t, esperado, resultado)
	mockRepo.AssertExpectations(t)
}
