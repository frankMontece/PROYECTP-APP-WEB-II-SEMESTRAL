package models

import "time"

type Obligacion struct {
	ID               uint       `json:"id" gorm:"primaryKey"`
	ResidenteID      uint       `json:"residente_id" gorm:"not null;index"`
	Residente        Usuario    `json:"residente,omitempty" gorm:"foreignKey:ResidenteID;references:ID"`
	Tipo             string     `json:"tipo" gorm:"size:20;not null"`
	Monto            float64    `json:"monto" gorm:"not null"`
	Periodo          string     `json:"periodo" gorm:"size:7;not null"`
	Estado           string     `json:"estado" gorm:"size:20;not null;default:pendiente"`
	FechaEmision     time.Time  `json:"fecha_emision"`
	FechaVencimiento time.Time  `json:"fecha_vencimiento"`
	FechaPago        *time.Time `json:"fecha_pago,omitempty"`
	Comprobante      string     `json:"comprobante,omitempty" gorm:"size:100"`
	MoraCalculada    float64    `json:"mora_calculada" gorm:"default:0"`
}

// Multa representa una sanción discreta aplicada a un residente
// BelongsTo: Usuario (residente)
type Multa struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	ResidenteID  uint       `json:"residente_id" gorm:"not null;index"`
	Residente    Usuario    `json:"residente,omitempty" gorm:"foreignKey:ResidenteID;references:ID"`
	Motivo       string     `json:"motivo" gorm:"size:100;not null"`
	Monto        float64    `json:"monto" gorm:"not null"`
	Estado       string     `json:"estado" gorm:"size:20;not null;default:pendiente"`
	FechaEmision time.Time  `json:"fecha_emision"`
	FechaPago    *time.Time `json:"fecha_pago,omitempty"`
}
