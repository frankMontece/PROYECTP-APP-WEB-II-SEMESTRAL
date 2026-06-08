package storage

import (
	"sync"

	"condominio-api/internal/models"
)

type MemoriaParqueo struct {
	vehiculos      []models.Vehiculo
	nextVehiculoID uint

	visitas      []models.VisitaVehiculo
	nextVisitaID uint

	accesos      []models.AccesoVehiculo
	nextAccesoID uint

	mu sync.Mutex
}

func NuevaMemoriaParqueo() *MemoriaParqueo {
	return &MemoriaParqueo{
		vehiculos:      []models.Vehiculo{},
		nextVehiculoID: 1,
		visitas:        []models.VisitaVehiculo{},
		nextVisitaID:   1,
		accesos:        []models.AccesoVehiculo{},
		nextAccesoID:   1,
	}
}

func (m *MemoriaParqueo) ListarVehiculos() []models.Vehiculo {
	m.mu.Lock()
	defer m.mu.Unlock()

	copia := make([]models.Vehiculo, len(m.vehiculos))
	copy(copia, m.vehiculos)
	return copia
}

func (m *MemoriaParqueo) BuscarVehiculoPorID(id uint) (models.Vehiculo, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, v := range m.vehiculos {
		if v.ID == id {
			return v, true
		}
	}
	return models.Vehiculo{}, false
}

func (m *MemoriaParqueo) CrearVehiculo(v models.Vehiculo) models.Vehiculo {
	m.mu.Lock()
	defer m.mu.Unlock()

	v.ID = m.nextVehiculoID
	m.nextVehiculoID++
	m.vehiculos = append(m.vehiculos, v)
	return v
}

func (m *MemoriaParqueo) ActualizarVehiculo(id uint, datos models.Vehiculo) (models.Vehiculo, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, v := range m.vehiculos {
		if v.ID == id {
			datos.ID = id
			m.vehiculos[i] = datos
			return datos, true
		}
	}
	return models.Vehiculo{}, false
}

func (m *MemoriaParqueo) BorrarVehiculo(id uint) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, v := range m.vehiculos {
		if v.ID == id {
			m.vehiculos = append(m.vehiculos[:i], m.vehiculos[i+1:]...)
			return true
		}
	}
	return false
}
