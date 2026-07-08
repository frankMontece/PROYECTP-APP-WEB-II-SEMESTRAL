package storage

import (
	"condominio-api/internal/models"
	"time"

	"gorm.io/gorm"
)

// VehiculoSQLite implementa VehiculosRepository usando GORM sobre SQLite
type VehiculoSQLite struct {
	db *gorm.DB
}

// NewVehiculoSQLite construye el repositorio de vehículos inyectándole la conexión ya abierta
func NewVehiculoSQLite(db *gorm.DB) *VehiculoSQLite {
	return &VehiculoSQLite{db: db}
}

// VEHICULOS

func (s *VehiculoSQLite) ListarVehiculos() []models.Vehiculo {
	var vehiculos []models.Vehiculo
	s.db.Find(&vehiculos)
	return vehiculos
}

func (s *VehiculoSQLite) BuscarVehiculoPorID(id uint) (models.Vehiculo, bool) {
	var v models.Vehiculo
	if err := s.db.First(&v, id).Error; err != nil {
		return models.Vehiculo{}, false
	}
	return v, true
}

func (s *VehiculoSQLite) CrearVehiculo(v models.Vehiculo) models.Vehiculo {
	v.CreatedAt = time.Now()
	s.db.Create(&v)
	return v
}

func (s *VehiculoSQLite) ActualizarVehiculo(id uint, datos models.Vehiculo) (models.Vehiculo, bool) {
	var existente models.Vehiculo
	if err := s.db.First(&existente, id).Error; err != nil {
		return models.Vehiculo{}, false
	}
	datos.ID = id // el ID no puede cambiar
	s.db.Save(&datos)
	return datos, true
}

func (s *VehiculoSQLite) BorrarVehiculo(id uint) bool {
	res := s.db.Delete(&models.Vehiculo{}, id)
	return res.RowsAffected > 0
}
