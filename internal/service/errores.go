package service

import "errors"

// Errores de dominio - centralizados para toda la aplicación
var (
	ErrNombreVacio            = errors.New("nombre es requerido")
	ErrCapacidadInvalida      = errors.New("capacidad debe ser mayor a 0")
	ErrFechasInvalidas        = errors.New("fecha_fin debe ser posterior a fecha_inicio")
	ErrNumeroPersonasInvalido = errors.New("numero_personas debe ser mayor a 0")
	ErrNoEncontrado           = errors.New("recurso no encontrado")
	ErrEmailEnUso             = errors.New("email ya en uso")
	ErrCredencialesInvalidas  = errors.New("credenciales inválidas")
	ErrSinPermiso             = errors.New("sin permisos para esta acción")
)
