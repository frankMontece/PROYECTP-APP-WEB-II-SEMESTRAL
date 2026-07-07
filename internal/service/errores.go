package service

import "errors"

// ============================================================================
// Errores generales
// ============================================================================

var (
	ErrNoEncontrado = errors.New("recurso no encontrado")
)

// ============================================================================
// Usuarios / Autenticación
// ============================================================================

var (
	ErrEmailEnUso            = errors.New("email ya en uso")
	ErrCredencialesInvalidas = errors.New("credenciales inválidas")
	ErrSinPermiso            = errors.New("sin permisos para esta acción")
)

// ============================================================================
// Áreas Sociales
// ============================================================================

var (
	ErrNombreVacio            = errors.New("nombre es requerido")
	ErrCapacidadInvalida      = errors.New("capacidad debe ser mayor a 0")
	ErrFechasInvalidas        = errors.New("fecha_fin debe ser posterior a fecha_inicio")
	ErrNumeroPersonasInvalido = errors.New("numero_personas debe ser mayor a 0")
)

// ============================================================================
// Obligaciones
// ============================================================================

var (
	ErrResidenteIDInvalido = errors.New("residente_id es requerido")
	ErrMontoInvalido       = errors.New("monto debe ser mayor a 0")
	ErrPeriodoVacio        = errors.New("periodo es requerido (formato YYYY-MM)")
	ErrMotivoVacio         = errors.New("motivo es requerido")
	ErrTipoObligacion      = errors.New("tipo debe ser 'mensual' o 'extraordinaria'")
)

// ============================================================================
// Parqueo
// ============================================================================

var (
// Agrega aquí los errores propios del módulo de parqueo.
// Ejemplo:
// ErrPlacaVacia        = errors.New("la placa es requerida")
// ErrVehiculoNoActivo  = errors.New("el vehículo no está activo")
)
