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
// Parqueo — Vehículos
// ============================================================================

var (
	ErrVehiculoNoEncontrado = errors.New("vehículo no encontrado")
	ErrPlacaRequerida       = errors.New("placa es requerida")
	ErrMarcaRequerida       = errors.New("marca es requerida")
	ErrResidenteRequerido   = errors.New("residente_id es requerido")
)

// ============================================================================
// Parqueo — Visitas
// ============================================================================

var (
	ErrVisitaNoEncontrada       = errors.New("visita no encontrada")
	ErrPlacaVisitanteRequerida  = errors.New("placa_visitante es requerida")
	ErrNombreVisitanteRequerido = errors.New("nombre_visitante es requerido")
	ErrCondominioRequerido      = errors.New("condominio_id es requerido")
)

// ============================================================================
// Parqueo — Accesos
// ============================================================================

var (
	ErrAccesoNoEncontrado     = errors.New("acceso no encontrado")
	ErrVehiculoRequerido      = errors.New("vehiculo_id es requerido")
	ErrTipoMovimientoInvalido = errors.New("tipo_movimiento debe ser 'entrada' o 'salida'")
)
