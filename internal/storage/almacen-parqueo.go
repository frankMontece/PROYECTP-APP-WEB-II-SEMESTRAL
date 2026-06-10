package storage

import "condominio-api/internal/models"

type AlmacenParqueo interface {
	//Vehiculos de residentes
	ListarVehiculos() []models.Vehiculo
	BuscarVehiculoPorID(id uint) (models.Vehiculo, bool)
	CrearVehiculo(v models.Vehiculo) models.Vehiculo
	ActualizarVehiculo(id uint, datos models.Vehiculo) (models.Vehiculo, bool)
	BorrarVehiculo(id uint) bool

	//Visitas de vehículos
	ListarVisitas() []models.VisitaVehiculo
	BuscarVisitaPorID(id uint) (models.VisitaVehiculo, bool)
	CrearVisita(v models.VisitaVehiculo) models.VisitaVehiculo
	ActualizarVisita(id uint, datos models.VisitaVehiculo) (models.VisitaVehiculo, bool)
	RegistrarSalidaVisita(id uint) (models.VisitaVehiculo, bool)
	BorrarVisita(id uint) bool

	//Bitácora de accesos
	ListarAccesos() []models.AccesoVehiculo
	BuscarAccesoPorID(id uint) (models.AccesoVehiculo, bool)
	CrearAcceso(a models.AccesoVehiculo) models.AccesoVehiculo
	BorrarAcceso(id uint) bool
}

var _ AlmacenParqueo = (*MemoriaParqueo)(nil)
