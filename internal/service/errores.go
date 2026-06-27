package service

import "errors"

// ─── Autenticación ────────────────────────────────────────────────────────────
var (
	ErrCredencialesInvalidas = errors.New("email o contraseña incorrectos")
	ErrEmailEnUso            = errors.New("el email ya está registrado")
)

// ─── Módulo A — Obligaciones (Héctor Fernández) ───────────────────────────────
var (
	ErrNoEncontrado        = errors.New("recurso no encontrado")
	ErrResidenteIDInvalido = errors.New("residente_id es requerido")
	ErrMontoInvalido       = errors.New("monto debe ser mayor a 0")
	ErrPeriodoVacio        = errors.New("periodo es requerido (formato YYYY-MM)")
	ErrMotivoVacio         = errors.New("motivo es requerido")
	ErrTipoObligacion      = errors.New("tipo debe ser 'mensual' o 'extraordinaria'")
)
