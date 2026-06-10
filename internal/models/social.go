package models

import "time"

// AreaSocial representa una área social
type AreaSocial struct {
	ID          uint   `json:"id"`
	Nombre      string `json:"nombre"` // requerido
	Descripcion string `json:"descripcion"`
	Capacidad   int    `json:"capacidad"` // requerido, mayor a 0
	Activo      bool   `json:"activo"`
}

// ReservaArea representa la reserva de un AreaSocial por un residente
type ReservaArea struct {
	ID             uint      `json:"id"`
	AreaID         uint      `json:"area_id"`         // requerido
	ResidenteID    uint      `json:"residente_id"`    // requerido
	FechaInicio    time.Time `json:"fecha_inicio"`    // requerido
	FechaFin       time.Time `json:"fecha_fin"`       // requerido, debe ser > FechaInicio
	Proposito      string    `json:"proposito"`       // requerido
	NumeroPersonas int       `json:"numero_personas"` // requerido, mayor a 0
	Estado         string    `json:"estado"`          // "pendiente" | "aprobada" | "cancelada"
	FechaCreacion  time.Time `json:"fecha_creacion"`
}

// Notificacion representa un mensaje interno enviado a un residente
type Notificacion struct {
	ID            uint      `json:"id"`
	ResidenteID   uint      `json:"residente_id"` // requerido
	Tipo          string    `json:"tipo"`         // requerido: "reserva", "pago", "aviso", etc.
	Mensaje       string    `json:"mensaje"`      // requerido
	Leida         bool      `json:"leida"`
	FechaCreacion time.Time `json:"fecha_creacion"`
}
