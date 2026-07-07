package models

import "time"

// AreaSocial representa un área social del condominio
type AreaSocial struct {
	ID          uint   `json:"id"          gorm:"primaryKey"`
	Nombre      string `json:"nombre"       gorm:"not null"`
	Descripcion string `json:"descripcion"`
	Capacidad   int    `json:"capacidad"    gorm:"not null"`
	Activo      bool   `json:"activo"       gorm:"default:true"`
}

// ReservaArea representa la reserva de un área social
type ReservaArea struct {
	ID             uint      `json:"id"              gorm:"primaryKey"`
	AreaID         uint      `json:"area_id"         gorm:"not null"`
	ResidenteID    uint      `json:"residente_id"    gorm:"not null"`
	FechaInicio    time.Time `json:"fecha_inicio"`
	FechaFin       time.Time `json:"fecha_fin"`
	Proposito      string    `json:"proposito"       gorm:"not null"`
	NumeroPersonas int       `json:"numero_personas" gorm:"not null"`
	Estado         string    `json:"estado"          gorm:"default:pendiente"`
	FechaCreacion  time.Time `json:"fecha_creacion"  gorm:"autoCreateTime"`
}

// Notificacion representa un mensaje interno del sistema
type Notificacion struct {
	ID            uint      `json:"id"             gorm:"primaryKey"`
	ResidenteID   uint      `json:"residente_id"   gorm:"not null"`
	Tipo          string    `json:"tipo"           gorm:"not null"`
	Mensaje       string    `json:"mensaje"        gorm:"not null"`
	Leida         bool      `json:"leida"          gorm:"default:false"`
	FechaCreacion time.Time `json:"fecha_creacion" gorm:"autoCreateTime"`
}
