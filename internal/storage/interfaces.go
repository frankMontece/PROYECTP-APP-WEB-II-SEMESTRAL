package storage

import "condominio-api/internal/models"

// ─── INTERFACES PEQUEÑAS (ISP) ──────────────────────────────────────────────

type AreaSocialRepository interface {
	ListarAreas() []models.AreaSocial
	BuscarAreaPorID(id uint) (models.AreaSocial, bool)
	CrearArea(a models.AreaSocial) models.AreaSocial
	ActualizarArea(id uint, datos models.AreaSocial) (models.AreaSocial, bool)
	BorrarArea(id uint) bool
}

type ReservaAreaRepository interface {
	ListarReservas() []models.ReservaArea
	BuscarReservaPorID(id uint) (models.ReservaArea, bool)
	CrearReserva(r models.ReservaArea) models.ReservaArea
	ActualizarReserva(id uint, datos models.ReservaArea) (models.ReservaArea, bool)
	BorrarReserva(id uint) bool
}

type NotificacionRepository interface {
	ListarNotificaciones() []models.Notificacion
	BuscarNotificacionPorID(id uint) (models.Notificacion, bool)
	CrearNotificacion(n models.Notificacion) models.Notificacion
	MarcarComoLeida(id uint) (models.Notificacion, bool)
	BorrarNotificacion(id uint) bool
}

type UserRepository interface {
	CrearUsuario(u models.Usuario) (models.Usuario, error)
	BuscarUsuarioPorEmail(email string) (models.Usuario, bool)
}
