package storage

import (
	"fmt"
	"time"

	"condominio-api/internal/models"

	"gorm.io/gorm"
)

func SembrarParqueo(db *gorm.DB) {
	var n int64
	db.Model(&models.Vehiculo{}).Count(&n)
	if n > 0 {
		return
	}

	ahora := time.Now()

	vehiculos := []models.Vehiculo{
		{ResidenteID: 1, Placa: "PBG-2241", Marca: "Toyota", Modelo: "Corolla", Color: "Blanco", Activo: true, CreatedAt: ahora},
		{ResidenteID: 2, Placa: "ACG-0873", Marca: "Chevrolet", Modelo: "Sail", Color: "Gris", Activo: true, CreatedAt: ahora},
		{ResidenteID: 3, Placa: "MBN-5519", Marca: "Kia", Modelo: "Sportage", Color: "Negro", Activo: true, CreatedAt: ahora},
	}
	db.Create(&vehiculos)

	entrada := ahora.Add(-2 * time.Hour)
	visitas := []models.VisitaVehiculo{
		{
			CondominioID:    1,
			ResidenteID:     1,
			PlacaVisitante:  "HBA-7731",
			NombreVisitante: "Ana Lucía Cedeño",
			Motivo:          "Visita familiar",
			CodigoQR:        fmt.Sprintf("QR-%d", ahora.UnixNano()),
			EstadoQR:        "pendiente",
			HoraEntrada:     &entrada,
		},
	}
	db.Create(&visitas)

	accesos := []models.AccesoVehiculo{
		{VehiculoID: 1, CondominioID: 1, TipoMovimiento: "entrada", FechaHora: ahora.Add(-3 * time.Hour), Observacion: "Sin novedad"},
		{VehiculoID: 2, CondominioID: 1, TipoMovimiento: "entrada", FechaHora: ahora.Add(-1 * time.Hour), Observacion: "Ingreso nocturno"},
	}
	db.Create(&accesos)

	fmt.Printf("✅ Datos de ejemplo sembrados en Parqueo: %d vehículos, %d visitas, %d accesos\n",
		len(vehiculos), len(visitas), len(accesos))
}
