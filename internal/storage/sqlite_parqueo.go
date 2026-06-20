package storage

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"condominio-api/internal/models"
)

// SQLiteParqueo implementa AlmacenParqueo usando GORM sobre SQLite.
// Es el reemplazo directo de MemoriaParqueo: misma interfaz, distinto backend.
// Para volver a memoria en main.go basta cambiar NewSQLiteParqueo por NuevaMemoriaParqueo.
type SQLiteParqueo struct {
	db *gorm.DB
}

// NewSQLiteParqueo construye el repositorio inyectándole la conexión ya abierta.
// La conexión se abre una sola vez en main.go y se comparte con UsuarioGORM.
func NewSQLiteParqueo(db *gorm.DB) *SQLiteParqueo {
	return &SQLiteParqueo{db: db}
}

// Verificación en compilación: si SQLiteParqueo no implementa todos los métodos
// de AlmacenParqueo, el compilador da error aquí antes de que el programa corra.
var _ AlmacenParqueo = (*SQLiteParqueo)(nil)

// SembrarVacio inserta datos de ejemplo solo si las tablas están vacías.
// Se llama una vez en main.go justo después de crear el almacén,
// igual que el profesor hace con almacen.SembrarVacio() en cafetería.
// Si ya hay datos, no hace nada: es seguro llamarlo en cada arranque.
func (s *SQLiteParqueo) SembrarVacio() {
	var n int64
	s.db.Model(&models.Vehiculo{}).Count(&n)
	if n > 0 {
		return
	}

	ahora := time.Now()

	vehiculos := []models.Vehiculo{
		{ResidenteID: 1, Placa: "PBG-2241", Marca: "Toyota", Modelo: "Corolla", Color: "Blanco", Activo: true, CreatedAt: ahora},
		{ResidenteID: 2, Placa: "ACG-0873", Marca: "Chevrolet", Modelo: "Sail", Color: "Gris", Activo: true, CreatedAt: ahora},
		{ResidenteID: 3, Placa: "MBN-5519", Marca: "Kia", Modelo: "Sportage", Color: "Negro", Activo: true, CreatedAt: ahora},
	}
	s.db.Create(&vehiculos)

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
	s.db.Create(&visitas)

	accesos := []models.AccesoVehiculo{
		{VehiculoID: 1, CondominioID: 1, TipoMovimiento: "entrada", FechaHora: ahora.Add(-3 * time.Hour), Observacion: "Sin novedad"},
		{VehiculoID: 2, CondominioID: 1, TipoMovimiento: "entrada", FechaHora: ahora.Add(-1 * time.Hour), Observacion: "Ingreso nocturno"},
	}
	s.db.Create(&accesos)
}

// =================================================================
// VEHICULOS
// =================================================================

func (s *SQLiteParqueo) ListarVehiculos() []models.Vehiculo {
	var vehiculos []models.Vehiculo
	s.db.Find(&vehiculos)
	return vehiculos
}

func (s *SQLiteParqueo) BuscarVehiculoPorID(id uint) (models.Vehiculo, bool) {
	var v models.Vehiculo
	if err := s.db.First(&v, id).Error; err != nil {
		return models.Vehiculo{}, false
	}
	return v, true
}

func (s *SQLiteParqueo) CrearVehiculo(v models.Vehiculo) models.Vehiculo {
	v.CreatedAt = time.Now()
	s.db.Create(&v)
	return v
}

func (s *SQLiteParqueo) ActualizarVehiculo(id uint, datos models.Vehiculo) (models.Vehiculo, bool) {
	var existente models.Vehiculo
	if err := s.db.First(&existente, id).Error; err != nil {
		return models.Vehiculo{}, false
	}
	datos.ID = id // el ID no puede cambiar
	s.db.Save(&datos)
	return datos, true
}

func (s *SQLiteParqueo) BorrarVehiculo(id uint) bool {
	res := s.db.Delete(&models.Vehiculo{}, id)
	return res.RowsAffected > 0
}

// =================================================================
// VISITAS
// =================================================================

func (s *SQLiteParqueo) ListarVisitas() []models.VisitaVehiculo {
	var visitas []models.VisitaVehiculo
	s.db.Find(&visitas)
	return visitas
}

func (s *SQLiteParqueo) BuscarVisitaPorID(id uint) (models.VisitaVehiculo, bool) {
	var vis models.VisitaVehiculo
	if err := s.db.First(&vis, id).Error; err != nil {
		return models.VisitaVehiculo{}, false
	}
	return vis, true
}

func (s *SQLiteParqueo) CrearVisita(vis models.VisitaVehiculo) models.VisitaVehiculo {
	s.db.Create(&vis)
	return vis
}

func (s *SQLiteParqueo) ActualizarVisita(id uint, datos models.VisitaVehiculo) (models.VisitaVehiculo, bool) {
	var existente models.VisitaVehiculo
	if err := s.db.First(&existente, id).Error; err != nil {
		return models.VisitaVehiculo{}, false
	}
	datos.ID = id
	s.db.Save(&datos)
	return datos, true
}

// RegistrarSalidaVisita marca la hora de salida y expira el QR.
// Usamos s.db.Save para que el cambio persista en disco:
// el próximo reinicio del servidor conservará el registro de salida.
func (s *SQLiteParqueo) RegistrarSalidaVisita(id uint) (models.VisitaVehiculo, bool) {
	var vis models.VisitaVehiculo
	if err := s.db.First(&vis, id).Error; err != nil {
		return models.VisitaVehiculo{}, false
	}
	ahora := time.Now()
	vis.HoraSalida = &ahora
	vis.EstadoQR = "expirado"
	s.db.Save(&vis)
	return vis, true
}

func (s *SQLiteParqueo) BorrarVisita(id uint) bool {
	res := s.db.Delete(&models.VisitaVehiculo{}, id)
	return res.RowsAffected > 0
}

// =================================================================
// ACCESOS
// =================================================================

func (s *SQLiteParqueo) ListarAccesos() []models.AccesoVehiculo {
	var accesos []models.AccesoVehiculo
	s.db.Find(&accesos)
	return accesos
}

func (s *SQLiteParqueo) BuscarAccesoPorID(id uint) (models.AccesoVehiculo, bool) {
	var a models.AccesoVehiculo
	if err := s.db.First(&a, id).Error; err != nil {
		return models.AccesoVehiculo{}, false
	}
	return a, true
}

func (s *SQLiteParqueo) CrearAcceso(a models.AccesoVehiculo) models.AccesoVehiculo {
	s.db.Create(&a)
	return a
}

func (s *SQLiteParqueo) BorrarAcceso(id uint) bool {
	res := s.db.Delete(&models.AccesoVehiculo{}, id)
	return res.RowsAffected > 0
}
