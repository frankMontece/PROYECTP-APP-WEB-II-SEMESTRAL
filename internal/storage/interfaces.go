package storage

import "condominio-api/internal/models"

// AreaSocialRepository define las operaciones para Áreas Sociales
type AreaSocialRepository interface {
	ListarAreas() []models.AreaSocial
	BuscarAreaPorID(id uint) (models.AreaSocial, bool)
	CrearArea(a models.AreaSocial) models.AreaSocial
	ActualizarArea(id uint, datos models.AreaSocial) (models.AreaSocial, bool)
	BorrarArea(id uint) bool
}

// ReservaAreaRepository define las operaciones para Reservas
type ReservaAreaRepository interface {
	ListarReservas() []models.ReservaArea
	BuscarReservaPorID(id uint) (models.ReservaArea, bool)
	CrearReserva(r models.ReservaArea) models.ReservaArea
	ActualizarReserva(id uint, datos models.ReservaArea) (models.ReservaArea, bool)
	BorrarReserva(id uint) bool
}

// NotificacionRepository define las operaciones para Notificaciones
type NotificacionRepository interface {
	ListarNotificaciones() []models.Notificacion
	BuscarNotificacionPorID(id uint) (models.Notificacion, bool)
	CrearNotificacion(n models.Notificacion) models.Notificacion
	MarcarComoLeida(id uint) (models.Notificacion, bool)
	BorrarNotificacion(id uint) bool
}

// UserRepository define las operaciones para Usuarios
type UserRepository interface {
	CrearUsuario(u models.Usuario) (models.Usuario, error)
	BuscarUsuarioPorEmail(email string) (models.Usuario, bool)
}

var _ AreaSocialRepository = (*AreaSQLite)(nil)
var _ ReservaAreaRepository = (*ReservaSQLite)(nil)
var _ NotificacionRepository = (*NotificacionSQLite)(nil)
var _ UserRepository = (*UsuarioGORM)(nil)
