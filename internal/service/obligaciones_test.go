package service

import (
	"testing"

	"condominio-api/internal/models"
)

type mockObligacionRepo struct {
	obligaciones []models.Obligacion
}

func (m *mockObligacionRepo) ListarObligaciones() []models.Obligacion {
	return m.obligaciones
}

func (m *mockObligacionRepo) BuscarObligacionPorID(id uint) (models.Obligacion, bool) {
	for _, o := range m.obligaciones {
		if o.ID == id {
			return o, true
		}
	}
	return models.Obligacion{}, false
}

func (m *mockObligacionRepo) CrearObligacion(o models.Obligacion) models.Obligacion {
	o.ID = uint(len(m.obligaciones) + 1)
	m.obligaciones = append(m.obligaciones, o)
	return o
}

func (m *mockObligacionRepo) ActualizarObligacion(id uint, datos models.Obligacion) (models.Obligacion, bool) {
	return models.Obligacion{}, false // no se usa en este test
}

func (m *mockObligacionRepo) BorrarObligacion(id uint) bool {
	return false // no se usa en este test
}

func TestCrearObligacion_RechazaMontoInvalido(t *testing.T) {
	mockRepo := &mockObligacionRepo{}
	svc := NewObligacionesService(mockRepo)

	_, err := svc.CrearObligacion(models.Obligacion{
		ResidenteID: 1,
		Tipo:        "mensual",
		Monto:       -50.00, // dato inválido a propósito
		Periodo:     "2026-06",
	})

	if err == nil {
		t.Fatal("esperaba error por monto inválido, pero CrearObligacion no falló")
	}
	if err != ErrMontoInvalido {
		t.Fatalf("esperaba ErrMontoInvalido, obtuve: %v", err)
	}
	if len(mockRepo.obligaciones) != 0 {
		t.Fatal("la obligación inválida llegó al repositorio — la validación no la detuvo")
	}
}
