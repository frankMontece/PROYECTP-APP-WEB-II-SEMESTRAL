package storage

import "condominio-api/internal/models"

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
