package service

import (
	"testing"

	"condominio-api/internal/models"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// ── Mock del repositorio ──────────────────────────────────────────────────────

type mockAreaRepo struct {
	mock.Mock
}

func (m *mockAreaRepo) ListarAreas() []models.AreaSocial {
	args := m.Called()
	return args.Get(0).([]models.AreaSocial)
}
func (m *mockAreaRepo) BuscarAreaPorID(id uint) (models.AreaSocial, bool) {
	args := m.Called(id)
	return args.Get(0).(models.AreaSocial), args.Bool(1)
}
func (m *mockAreaRepo) CrearArea(a models.AreaSocial) models.AreaSocial {
	args := m.Called(a)
	return args.Get(0).(models.AreaSocial)
}
func (m *mockAreaRepo) ActualizarArea(id uint, datos models.AreaSocial) (models.AreaSocial, bool) {
	args := m.Called(id, datos)
	return args.Get(0).(models.AreaSocial), args.Bool(1)
}
func (m *mockAreaRepo) BorrarArea(id uint) bool {
	args := m.Called(id)
	return args.Bool(0)
}

// ── Tests ─────────────────────────────────────────────────────────────────────

func TestAreaService_CrearArea_NombreVacio(t *testing.T) {
	// Preparar
	repo := new(mockAreaRepo)
	svc := NewAreaService(repo)

	entrada := models.AreaSocial{Nombre: "", Capacidad: 50}

	// Ejecutar
	_, err := svc.CrearArea(entrada)

	// Verificar
	require.ErrorIs(t, err, ErrNombreVacio)
	repo.AssertNotCalled(t, "CrearArea") // la regla de negocio lo bloqueó
}

func TestAreaService_CrearArea_CapacidadInvalida(t *testing.T) {
	// Preparar
	repo := new(mockAreaRepo)
	svc := NewAreaService(repo)

	entrada := models.AreaSocial{Nombre: "Salón", Capacidad: 0}

	// Ejecutar
	_, err := svc.CrearArea(entrada)

	// Verificar
	require.ErrorIs(t, err, ErrCapacidadInvalida)
	repo.AssertNotCalled(t, "CrearArea")
}
