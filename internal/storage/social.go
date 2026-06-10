package storage

import (
	"sync"
	"time"

	"condominio-api/internal/models"
)

// AlmacenSocial define la interfaz que deben cumplir las implementaciones
// de almacenamiento (memoria, SQLite, etc.). Facilita los tests y el reemplazo futuro.
type AlmacenSocial interface {
	// Áreas Sociales
	ListarAreas() []models.AreaSocial
	BuscarAreaPorID(id uint) (models.AreaSocial, bool)
	CrearArea(a models.AreaSocial) models.AreaSocial
	ActualizarArea(id uint, datos models.AreaSocial) (models.AreaSocial, bool)
	BorrarArea(id uint) bool

	// Reservas
	ListarReservas() []models.ReservaArea
	BuscarReservaPorID(id uint) (models.ReservaArea, bool)
	CrearReserva(r models.ReservaArea) models.ReservaArea
	ActualizarReserva(id uint, datos models.ReservaArea) (models.ReservaArea, bool)
	BorrarReserva(id uint) bool

	// Notificaciones
	ListarNotificaciones() []models.Notificacion
	BuscarNotificacionPorID(id uint) (models.Notificacion, bool)
	CrearNotificacion(n models.Notificacion) models.Notificacion
	MarcarComoLeida(id uint) (models.Notificacion, bool)
	BorrarNotificacion(id uint) bool
}

// MemoriaSocial es la implementación en RAM de AlmacenSocial
type MemoriaSocial struct {
	mu sync.Mutex

	areas      []models.AreaSocial
	nextAreaID uint

	reservas      []models.ReservaArea
	nextReservaID uint

	notificaciones     []models.Notificacion
	nextNotificacionID uint
}

// NuevaMemoriaSocial inicializa el almacén con slices vacíos e IDs desde 1
func NuevaMemoriaSocial() *MemoriaSocial {
	return &MemoriaSocial{
		areas:              []models.AreaSocial{},
		nextAreaID:         1,
		reservas:           []models.ReservaArea{},
		nextReservaID:      1,
		notificaciones:     []models.Notificacion{},
		nextNotificacionID: 1,
	}
}

// ─────────────────────────────────────────────────────────
// ÁREAS SOCIALES
// ─────────────────────────────────────────────────────────

func (m *MemoriaSocial) ListarAreas() []models.AreaSocial {
	m.mu.Lock()
	defer m.mu.Unlock()
	copia := make([]models.AreaSocial, len(m.areas))
	copy(copia, m.areas)
	return copia
}

func (m *MemoriaSocial) BuscarAreaPorID(id uint) (models.AreaSocial, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, a := range m.areas {
		if a.ID == id {
			return a, true
		}
	}
	return models.AreaSocial{}, false
}

func (m *MemoriaSocial) CrearArea(a models.AreaSocial) models.AreaSocial {
	m.mu.Lock()
	defer m.mu.Unlock()
	a.ID = m.nextAreaID
	m.nextAreaID++
	m.areas = append(m.areas, a)
	return a
}

func (m *MemoriaSocial) ActualizarArea(id uint, datos models.AreaSocial) (models.AreaSocial, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, a := range m.areas {
		if a.ID == id {
			datos.ID = id
			m.areas[i] = datos
			return datos, true
		}
	}
	return models.AreaSocial{}, false
}

func (m *MemoriaSocial) BorrarArea(id uint) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, a := range m.areas {
		if a.ID == id {
			m.areas = append(m.areas[:i], m.areas[i+1:]...)
			return true
		}
	}
	return false
}

// ─────────────────────────────────────────────────────────
// RESERVAS DE ÁREAS
// ─────────────────────────────────────────────────────────

func (m *MemoriaSocial) ListarReservas() []models.ReservaArea {
	m.mu.Lock()
	defer m.mu.Unlock()
	copia := make([]models.ReservaArea, len(m.reservas))
	copy(copia, m.reservas)
	return copia
}

func (m *MemoriaSocial) BuscarReservaPorID(id uint) (models.ReservaArea, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, r := range m.reservas {
		if r.ID == id {
			return r, true
		}
	}
	return models.ReservaArea{}, false
}

func (m *MemoriaSocial) CrearReserva(r models.ReservaArea) models.ReservaArea {
	m.mu.Lock()
	defer m.mu.Unlock()
	r.ID = m.nextReservaID
	m.nextReservaID++
	r.FechaCreacion = time.Now()
	r.Estado = "pendiente"
	m.reservas = append(m.reservas, r)
	return r
}

func (m *MemoriaSocial) ActualizarReserva(id uint, datos models.ReservaArea) (models.ReservaArea, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, r := range m.reservas {
		if r.ID == id {
			datos.ID = id
			datos.FechaCreacion = r.FechaCreacion // preservar fecha original
			m.reservas[i] = datos
			return datos, true
		}
	}
	return models.ReservaArea{}, false
}

func (m *MemoriaSocial) BorrarReserva(id uint) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, r := range m.reservas {
		if r.ID == id {
			m.reservas = append(m.reservas[:i], m.reservas[i+1:]...)
			return true
		}
	}
	return false
}

// ─────────────────────────────────────────────────────────
// NOTIFICACIONES
// ─────────────────────────────────────────────────────────

func (m *MemoriaSocial) ListarNotificaciones() []models.Notificacion {
	m.mu.Lock()
	defer m.mu.Unlock()
	copia := make([]models.Notificacion, len(m.notificaciones))
	copy(copia, m.notificaciones)
	return copia
}

func (m *MemoriaSocial) BuscarNotificacionPorID(id uint) (models.Notificacion, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, n := range m.notificaciones {
		if n.ID == id {
			return n, true
		}
	}
	return models.Notificacion{}, false
}

func (m *MemoriaSocial) CrearNotificacion(n models.Notificacion) models.Notificacion {
	m.mu.Lock()
	defer m.mu.Unlock()
	n.ID = m.nextNotificacionID
	m.nextNotificacionID++
	n.Leida = false
	n.FechaCreacion = time.Now()
	m.notificaciones = append(m.notificaciones, n)
	return n
}

// MarcarComoLeida es una acción de negocio: marca leida=true sin reemplazar el objeto completo
func (m *MemoriaSocial) MarcarComoLeida(id uint) (models.Notificacion, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, n := range m.notificaciones {
		if n.ID == id {
			m.notificaciones[i].Leida = true
			return m.notificaciones[i], true
		}
	}
	return models.Notificacion{}, false
}

func (m *MemoriaSocial) BorrarNotificacion(id uint) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, n := range m.notificaciones {
		if n.ID == id {
			m.notificaciones = append(m.notificaciones[:i], m.notificaciones[i+1:]...)
			return true
		}
	}
	return false
}

// Asegurar que MemoriaSocial implementa AlmacenSocial en tiempo de compilación
var _ AlmacenSocial = (*MemoriaSocial)(nil)
