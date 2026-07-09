package models

import "time"

// Usuario representa un usuario del sistema
// Modelo DEFINITIVO - NO AGREGAR CAMPOS ADICIONALES
type Usuario struct {
	ID           int       `json:"id"         gorm:"primaryKey"`
	Rol          string    `json:"rol"        gorm:"not null"`
	Email        string    `json:"email"      gorm:"unique;not null"`
	PasswordHash string    `json:"password"   gorm:"not null"`
	CreadoEn     time.Time `json:"creado_en"  gorm:"autoCreateTime"`
}
