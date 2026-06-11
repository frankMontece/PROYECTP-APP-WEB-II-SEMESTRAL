package storage

import (
	"condominio-api/internal/models"
	"fmt"
	"sync"
)

type MemoriaObligaciones struct {
	Obligaciones     []models.Obligacion
	nextObligacionID uint
	Multas           []models.Multa
	nextMultaID      uint
	mu               sync.Mutex
}

func NuevaMemoriaObligaciones() *MemoriaObligaciones {
	return &MemoriaObligaciones{
		Obligaciones:     []models.Obligacion{},
		nextObligacionID: 1,
		Multas:           []models.Multa{},
		nextMultaID:      1,
	}
}

// =========================================================
// OBLIGACIONES
// =========================================================

func (m *MemoriaObligaciones) ListarObligaciones() ([]models.Obligacion, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	copia := make([]models.Obligacion, len(m.Obligaciones))
	copy(copia, m.Obligaciones)
	return copia, nil
}

func (m *MemoriaObligaciones) BuscarObligacionPorID(id uint) (*models.Obligacion, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i := range m.Obligaciones {
		if m.Obligaciones[i].ID == id {
			copia := m.Obligaciones[i] // copia segura
			return &copia, nil
		}
	}
	return nil, fmt.Errorf("obligacion con ID %d no encontrada", id)
}

func (m *MemoriaObligaciones) CrearObligacion(obligacion models.Obligacion) (*models.Obligacion, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	obligacion.ID = m.nextObligacionID
	m.nextObligacionID++
	m.Obligaciones = append(m.Obligaciones, obligacion)
	ultimo := len(m.Obligaciones) - 1
	return &m.Obligaciones[ultimo], nil
}

func (m *MemoriaObligaciones) ActualizarObligacion(id uint, obligacion models.Obligacion) (*models.Obligacion, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i := range m.Obligaciones {
		if m.Obligaciones[i].ID == id {
			obligacion.ID = id // preserva el ID original
			m.Obligaciones[i] = obligacion
			return &m.Obligaciones[i], nil
		}
	}
	return nil, fmt.Errorf("obligacion con ID %d no encontrada", id)
}

func (m *MemoriaObligaciones) EliminarObligacion(id uint) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i := range m.Obligaciones {
		if m.Obligaciones[i].ID == id {
			m.Obligaciones = append(m.Obligaciones[:i], m.Obligaciones[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("obligacion con ID %d no encontrada", id)
}

// =========================================================
// MULTAS
// =========================================================
func (m *MemoriaObligaciones) ListarMultas() ([]models.Multa, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	copia := make([]models.Multa, len(m.Multas))
	copy(copia, m.Multas)
	return copia, nil
}

func (m *MemoriaObligaciones) BuscarMultaPorID(id uint) (*models.Multa, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i := range m.Multas {
		if m.Multas[i].ID == id {
			copia := m.Multas[i] // copia segura
			return &copia, nil
		}
	}
	return nil, fmt.Errorf("multa con ID %d no encontrada", id)
}

func (m *MemoriaObligaciones) CrearMulta(multa models.Multa) (*models.Multa, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	multa.ID = m.nextMultaID
	m.nextMultaID++
	m.Multas = append(m.Multas, multa)
	ultimo := len(m.Multas) - 1
	return &m.Multas[ultimo], nil
}

func (m *MemoriaObligaciones) ActualizarMulta(id uint, multa models.Multa) (*models.Multa, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i := range m.Multas {
		if m.Multas[i].ID == id {
			multa.ID = id // preserva el ID original
			m.Multas[i] = multa
			return &m.Multas[i], nil
		}
	}
	return nil, fmt.Errorf("multa con ID %d no encontrada", id)
}

func (m *MemoriaObligaciones) EliminarMulta(id uint) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i := range m.Multas {
		if m.Multas[i].ID == id {
			m.Multas = append(m.Multas[:i], m.Multas[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("multa con ID %d no encontrada", id)
}
