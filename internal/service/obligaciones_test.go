package service

import (
	"testing"

	"condominio-api/internal/models"
)

// mockObligacionRepo es un doble de prueba — implementa ObligacionRepository
// pero no toca ninguna base de datos real. Solo guarda en un slice en RAM
// para poder verificar si CrearObligacion llegó a llamarlo o no.
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

// TestCrearObligacion_RechazaMontoInvalido prueba una regla de negocio real:
// un monto <= 0 NUNCA debe llegar al repositorio.
// Qué se rompería si fallara: si alguien borra el "if o.Monto <= 0" del
// service, este test falla porque el mock terminaría con 1 elemento
// guardado en lugar de 0 — la validación dejaría de proteger los datos.
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
