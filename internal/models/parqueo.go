package models

import "time"

type Vehiculo struct {
	ID          uint      `json:"id"`
	ResidenteID uint      `json:"residente_id"`
	Placa       string    `json:"placa"`
	Marca       string    `json:"marca"`
	Modelo      string    `json:"modelo"`
	Color       string    `json:"color"`
	Activo      bool      `json:"activo"`
	CreatedAt   time.Time `json:"created_at"`
}

type VisitaVehiculo struct {
	ID              uint       `json:"id"`
	CondominioID    uint       `json:"condominio_id"`
	ResidenteID     uint       `json:"residente_id"`
	PlacaVisitante  string     `json:"placa_visitante"`
	NombreVisitante string     `json:"nombre_visitante"`
	HoraEntrada     *time.Time `json:"hora_entrada"`
	HoraSalida      *time.Time `json:"hora_salida"`
	Motivo          string     `json:"motivo"`
	CodigoQR        string     `json:"codigo_qr"`
	EstadoQR        string     `json:"estado_qr"`
}

type AccesoVehiculo struct {
	ID             uint      `json:"id"`
	VehiculoID     uint      `json:"vehiculo_id"`
	CondominioID   uint      `json:"condominio_id"`
	TipoMovimiento string    `json:"tipo_movimiento"`
	FechaHora      time.Time `json:"fecha_hora"`
	Observacion    string    `json:"observacion"`
}
