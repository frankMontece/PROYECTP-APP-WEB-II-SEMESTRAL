package models

import "time"

type Obligacion struct {
	ID               uint       `json:"id"` // Identificador único de la obligación, uint es para evitar problemas con IDs negativos
	ResidenteID      uint       `json:"residente_id"`
	Tipo             string     `json:"tipo"`    // "mensual", "extraordinaria"
	Monto            float64    `json:"monto"`   // Monto de la obligación
	Periodo          string     `json:"periodo"` // Formato "YYYY-MM"
	Estado           string     `json:"estado"`  // "pendiente", "pagada", "vencida"
	FechaEmision     time.Time  `json:"fecha_emision"`
	FechaVencimiento time.Time  `json:"fecha_vencimiento"`
	FechaPago        *time.Time `json:"fecha_pago,omitempty"` // Fecha de pago, nula si no se ha pagado, se usa puntero para permitir nulo
	Comprobante      string     `json:"comprobante,omitempty"`
	MoraCalculada    float64    `json:"mora_calculada"`
}

// Multa representa una sanción discreta aplicada a un residente
// BelongsTo: Residente
type Multa struct {
	ID           uint       `json:"id"`
	ResidenteID  uint       `json:"residente_id"`
	Motivo       string     `json:"motivo"` //
	Monto        float64    `json:"monto"`
	Estado       string     `json:"estado"` // "pendiente", "pagada", "apelada"
	FechaEmision time.Time  `json:"fecha_emision"`
	FechaPago    *time.Time `json:"fecha_pago,omitempty"` // Fecha de pago, nula si no se ha pagado, se usa puntero para permitir nulo
}
