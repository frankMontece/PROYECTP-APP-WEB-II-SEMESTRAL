package models

import "time"

type Vehiculo struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	ResidenteID uint      `json:"residente_id" gorm:"not null;index"`
	Placa       string    `json:"placa" gorm:"unique;not null"`
	Marca       string    `json:"marca" gorm:"not null"`
	Modelo      string    `json:"modelo"`
	Color       string    `json:"color"`
	Activo      bool      `json:"activo: true"`
	CreatedAt   time.Time `json:"created_at"`
}

type VisitaVehiculo struct {
	ID              uint       `json:"id" gorm:"primaryKey"`
	CondominioID    uint       `json:"condominio_id" gorm:"not null;index"`
	ResidenteID     uint       `json:"residente_id" gorm:"not null;index"`
	PlacaVisitante  string     `json:"placa_visitante" gorm:"not null"`
	NombreVisitante string     `json:"nombre_visitante" gorm:"not null"`
	HoraEntrada     *time.Time `json:"hora_entrada"`
	HoraSalida      *time.Time `json:"hora_salida"`
	Motivo          string     `json:"motivo"`
	CodigoQR        string     `json:"codigo_qr" gorm:"unique"`
	EstadoQR        string     `json:"estado_qr"`
}

type AccesoVehiculo struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	VehiculoID     uint      `json:"vehiculo_id" gorm:"not null;index"`
	CondominioID   uint      `json:"condominio_id" gorm:"not null;index"`
	TipoMovimiento string    `json:"tipo_movimiento" gorm:"not null"`
	FechaHora      time.Time `json:"fecha_hora"`
	Observacion    string    `json:"observacion"`
}
