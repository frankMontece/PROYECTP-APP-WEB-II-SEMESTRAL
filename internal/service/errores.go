package service

import "errors"

// Errores de autenticación
var (
	ErrCredencialesInvalidas = errors.New("email o contraseña incorrectos")
	ErrEmailEnUso            = errors.New("el email ya está registrado")
)

// Errores de vehículos
var (
	ErrVehiculoNoEncontrado = errors.New("vehículo no encontrado")
	ErrPlacaRequerida       = errors.New("placa es requerida")
	ErrMarcaRequerida       = errors.New("marca es requerida")
	ErrResidenteRequerido   = errors.New("residente_id es requerido")
)

// Errores de visitas
var (
	ErrVisitaNoEncontrada       = errors.New("visita no encontrada")
	ErrPlacaVisitanteRequerida  = errors.New("placa_visitante es requerida")
	ErrNombreVisitanteRequerido = errors.New("nombre_visitante es requerido")
	ErrCondominioRequerido      = errors.New("condominio_id es requerido")
)

// Errores de accesos
var (
	ErrAccesoNoEncontrado     = errors.New("acceso no encontrado")
	ErrVehiculoRequerido      = errors.New("vehiculo_id es requerido")
	ErrTipoMovimientoInvalido = errors.New("tipo_movimiento debe ser 'entrada' o 'salida'")
)
