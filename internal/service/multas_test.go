package service

import (
	"testing"

	"condominio-api/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockMultaRepo implementa storage.MultaRepository usando testify/mock.
type mockMultaRepo struct {
	mock.Mock
}

func (m *mockMultaRepo) ListarMultas() []models.Multa {
	args := m.Called()
	return args.Get(0).([]models.Multa)
}

func (m *mockMultaRepo) BuscarMultaPorID(id uint) (models.Multa, bool) {
	args := m.Called(id)
	return args.Get(0).(models.Multa), args.Bool(1)
}

func (m *mockMultaRepo) CrearMulta(mu models.Multa) models.Multa {
	args := m.Called(mu)
	return args.Get(0).(models.Multa)
}

func (m *mockMultaRepo) ActualizarMulta(id uint, datos models.Multa) (models.Multa, bool) {
	args := m.Called(id, datos)
	return args.Get(0).(models.Multa), args.Bool(1)
}

func (m *mockMultaRepo) BorrarMulta(id uint) bool {
	args := m.Called(id)
	return args.Bool(0)
}

// ---------- CrearMulta ----------

func TestCrearMulta_RechazaResidenteIDInvalido(t *testing.T) {
	mockRepo := new(mockMultaRepo)
	svc := NewMultasService(mockRepo)

	_, err := svc.CrearMulta(models.Multa{
		ResidenteID: 0, // inválido
		Motivo:      "Ruido excesivo",
		Monto:       30.00,
	})

	assert.Equal(t, ErrResidenteIDInvalido, err)
	mockRepo.AssertNotCalled(t, "CrearMulta")
}

func TestCrearMulta_RechazaMotivoVacio(t *testing.T) {
	mockRepo := new(mockMultaRepo)
	svc := NewMultasService(mockRepo)

	_, err := svc.CrearMulta(models.Multa{
		ResidenteID: 1,
		Motivo:      "   ", // inválido
		Monto:       30.00,
	})

	assert.Equal(t, ErrMotivoVacio, err)
	mockRepo.AssertNotCalled(t, "CrearMulta")
}

func TestCrearMulta_RechazaMontoInvalido(t *testing.T) {
	mockRepo := new(mockMultaRepo)
	svc := NewMultasService(mockRepo)

	_, err := svc.CrearMulta(models.Multa{
		ResidenteID: 1,
		Motivo:      "Ruido excesivo",
		Monto:       0, // inválido
	})

	assert.Equal(t, ErrMontoInvalido, err)
	mockRepo.AssertNotCalled(t, "CrearMulta")
}

func TestCrearMulta_Exitosa(t *testing.T) {
	mockRepo := new(mockMultaRepo)
	svc := NewMultasService(mockRepo)

	entrada := models.Multa{
		ResidenteID: 1,
		Motivo:      "Ruido excesivo",
		Monto:       30.00,
	}
	entradaConEstado := entrada
	entradaConEstado.Estado = "pendiente"

	esperado := entradaConEstado
	esperado.ID = 1

	mockRepo.On("CrearMulta", entradaConEstado).Return(esperado)

	resultado, err := svc.CrearMulta(entrada)

	assert.NoError(t, err)
	assert.Equal(t, "pendiente", resultado.Estado)
	assert.Equal(t, uint(1), resultado.ID)
	mockRepo.AssertExpectations(t)
}

// ---------- ObtenerMulta ----------

func TestObtenerMulta_Encontrada(t *testing.T) {
	mockRepo := new(mockMultaRepo)
	svc := NewMultasService(mockRepo)

	esperado := models.Multa{ID: 5, ResidenteID: 1, Motivo: "Mascota sin correa", Monto: 15.00}
	mockRepo.On("BuscarMultaPorID", uint(5)).Return(esperado, true)

	resultado, err := svc.ObtenerMulta(5)

	assert.NoError(t, err)
	assert.Equal(t, esperado, resultado)
	mockRepo.AssertExpectations(t)
}

func TestObtenerMulta_NoEncontrada(t *testing.T) {
	mockRepo := new(mockMultaRepo)
	svc := NewMultasService(mockRepo)

	mockRepo.On("BuscarMultaPorID", uint(99)).Return(models.Multa{}, false)

	_, err := svc.ObtenerMulta(99)

	assert.Equal(t, ErrNoEncontrado, err)
	mockRepo.AssertExpectations(t)
}

// ---------- ActualizarMulta ----------

func TestActualizarMulta_RechazaResidenteIDInvalido(t *testing.T) {
	mockRepo := new(mockMultaRepo)
	svc := NewMultasService(mockRepo)

	_, err := svc.ActualizarMulta(1, models.Multa{
		ResidenteID: 0, // inválido
		Motivo:      "Ruido excesivo",
		Monto:       30.00,
	})

	assert.Equal(t, ErrResidenteIDInvalido, err)
	mockRepo.AssertNotCalled(t, "ActualizarMulta")
}

func TestActualizarMulta_RechazaMotivoVacio(t *testing.T) {
	mockRepo := new(mockMultaRepo)
	svc := NewMultasService(mockRepo)

	_, err := svc.ActualizarMulta(1, models.Multa{
		ResidenteID: 1,
		Motivo:      "",
		Monto:       30.00,
	})

	assert.Equal(t, ErrMotivoVacio, err)
	mockRepo.AssertNotCalled(t, "ActualizarMulta")
}

func TestActualizarMulta_RechazaMontoInvalido(t *testing.T) {
	mockRepo := new(mockMultaRepo)
	svc := NewMultasService(mockRepo)

	_, err := svc.ActualizarMulta(1, models.Multa{
		ResidenteID: 1,
		Motivo:      "Ruido excesivo",
		Monto:       -5,
	})

	assert.Equal(t, ErrMontoInvalido, err)
	mockRepo.AssertNotCalled(t, "ActualizarMulta")
}

func TestActualizarMulta_NoEncontrada(t *testing.T) {
	mockRepo := new(mockMultaRepo)
	svc := NewMultasService(mockRepo)

	datos := models.Multa{ResidenteID: 1, Motivo: "Ruido excesivo", Monto: 30.00}
	mockRepo.On("ActualizarMulta", uint(1), datos).Return(models.Multa{}, false)

	_, err := svc.ActualizarMulta(1, datos)

	assert.Equal(t, ErrNoEncontrado, err)
	mockRepo.AssertExpectations(t)
}

func TestActualizarMulta_Exitosa(t *testing.T) {
	mockRepo := new(mockMultaRepo)
	svc := NewMultasService(mockRepo)

	datos := models.Multa{ResidenteID: 1, Motivo: "Ruido excesivo", Monto: 30.00}
	esperado := datos
	esperado.ID = 1

	mockRepo.On("ActualizarMulta", uint(1), datos).Return(esperado, true)

	resultado, err := svc.ActualizarMulta(1, datos)

	assert.NoError(t, err)
	assert.Equal(t, esperado, resultado)
	mockRepo.AssertExpectations(t)
}

// ---------- BorrarMulta ----------

func TestBorrarMulta_Exitosa(t *testing.T) {
	mockRepo := new(mockMultaRepo)
	svc := NewMultasService(mockRepo)

	mockRepo.On("BorrarMulta", uint(1)).Return(true)

	err := svc.BorrarMulta(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestBorrarMulta_NoEncontrada(t *testing.T) {
	mockRepo := new(mockMultaRepo)
	svc := NewMultasService(mockRepo)

	mockRepo.On("BorrarMulta", uint(99)).Return(false)

	err := svc.BorrarMulta(99)

	assert.Equal(t, ErrNoEncontrado, err)
	mockRepo.AssertExpectations(t)
}

// ---------- ListarMultas ----------

func TestListarMultas(t *testing.T) {
	mockRepo := new(mockMultaRepo)
	svc := NewMultasService(mockRepo)

	esperado := []models.Multa{{ID: 1}, {ID: 2}}
	mockRepo.On("ListarMultas").Return(esperado)

	resultado := svc.ListarMultas()

	assert.Equal(t, esperado, resultado)
	mockRepo.AssertExpectations(t)
}
