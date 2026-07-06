package service

import (
	"testing"

	"condominio-api/internal/models"
)

// El test que prueba que esas reglas de validación realmente funcionan — que un monto negativo se rechaza antes de guardarse
type mockObligacionRepo struct { //se crea la BD falsa para probar los test
	obligaciones []models.Obligacion
}

func (m *mockObligacionRepo) ListarObligaciones() []models.Obligacion {
	return m.obligaciones
}

func (m *mockObligacionRepo) BuscarObligacionPorID(id uint) (models.Obligacion, bool) {
	for _, o := range m.obligaciones {
		if o.ID == id {
			return o, true //no se usa en este test, pero se implementa para cumplir con la interfaz
		}
	}
	return models.Obligacion{}, false
}

func (m *mockObligacionRepo) CrearObligacion(o models.Obligacion) models.Obligacion {
	o.ID = uint(len(m.obligaciones) + 1)
	m.obligaciones = append(m.obligaciones, o) //hay que listar todos los test pero solamente importa el de crearobligacion
	return o
}

func (m *mockObligacionRepo) ActualizarObligacion(id uint, datos models.Obligacion) (models.Obligacion, bool) {
	return models.Obligacion{}, false // no se usa en este test
}

func (m *mockObligacionRepo) BorrarObligacion(id uint) bool {
	return false // no se usa en este test
}

func TestCrearObligacion_RechazaMontoInvalido(t *testing.T) {
	mockRepo := &mockObligacionRepo{} //se crea el repo falso y se lo pasa al service para probar el test
	svc := NewObligacionesService(mockRepo)

	_, err := svc.CrearObligacion(models.Obligacion{ //se pide crear una obligación con un monto negativo
		ResidenteID: 1,
		Tipo:        "mensual",
		Monto:       -50.00, // dato inválido a propósito
		Periodo:     "2026-06",
	})

	if err == nil {
		t.Fatal("esperaba error por monto inválido, pero CrearObligacion no falló") //muestra que la validación no funcionó si no hay error
	}
	if err != ErrPeriodoVacio {
		t.Fatalf("esperaba ErrPeriodoVacio, obtuve: %v", err)
	}
	if len(mockRepo.obligaciones) != 0 {
		t.Fatal("la obligación inválida llegó al repositorio — la validación no la detuvo") //confirma que la validación cortó el flujo antes de llamar al repositorio.
	}
}
