package storage

import (
	"sync"
	"time"

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

// VISITAS

func (m *MemoriaParqueo) ListarVisitas() []models.VisitaVehiculo {
	m.mu.Lock()
	defer m.mu.Unlock()

	copia := make([]models.VisitaVehiculo, len(m.visitas))
	copy(copia, m.visitas)
	return copia
}

func (m *MemoriaParqueo) BuscarVisitaPorID(id uint) (models.VisitaVehiculo, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, vis := range m.visitas {
		if vis.ID == id {
			return vis, true
		}
	}
	return models.VisitaVehiculo{}, false
}

func (m *MemoriaParqueo) CrearVisita(vis models.VisitaVehiculo) models.VisitaVehiculo {
	m.mu.Lock()
	defer m.mu.Unlock()

	vis.ID = m.nextVisitaID
	m.nextVisitaID++
	m.visitas = append(m.visitas, vis)
	return vis
}

func (m *MemoriaParqueo) ActualizarVisita(id uint, datos models.VisitaVehiculo) (models.VisitaVehiculo, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, vis := range m.visitas {
		if vis.ID == id {
			datos.ID = id
			m.visitas[i] = datos
			return datos, true
		}
	}
	return models.VisitaVehiculo{}, false
}

func (m *MemoriaParqueo) RegistrarSalidaVisita(id uint) (models.VisitaVehiculo, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, vis := range m.visitas {
		if vis.ID == id {
			ahora := time.Now()
			m.visitas[i].HoraSalida = &ahora
			m.visitas[i].EstadoQR = "expirado"
			return m.visitas[i], true
		}
	}
	return models.VisitaVehiculo{}, false
}

func (m *MemoriaParqueo) BorrarVisita(id uint) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, vis := range m.visitas {
		if vis.ID == id {
			m.visitas = append(m.visitas[:i], m.visitas[i+1:]...)
			return true
		}
	}
	return false
}

// ACCESOS

func (m *MemoriaParqueo) ListarAccesos() []models.AccesoVehiculo {
	m.mu.Lock()
	defer m.mu.Unlock()

	copia := make([]models.AccesoVehiculo, len(m.accesos))
	copy(copia, m.accesos)
	return copia
}

func (m *MemoriaParqueo) BuscarAccesoPorID(id uint) (models.AccesoVehiculo, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, a := range m.accesos {
		if a.ID == id {
			return a, true
		}
	}
	return models.AccesoVehiculo{}, false
}

func (m *MemoriaParqueo) CrearAcceso(a models.AccesoVehiculo) models.AccesoVehiculo {
	m.mu.Lock()
	defer m.mu.Unlock()

	a.ID = m.nextAccesoID
	m.nextAccesoID++
	m.accesos = append(m.accesos, a)
	return a
}

func (m *MemoriaParqueo) BorrarAcceso(id uint) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, a := range m.accesos {
		if a.ID == id {
			m.accesos = append(m.accesos[:i], m.accesos[i+1:]...)
			return true
		}
	}
	return false
}
