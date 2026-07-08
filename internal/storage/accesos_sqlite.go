package storage

import (
	"condominio-api/internal/models"

	"gorm.io/gorm"
)

// AccesoSQLite implementa AccesosRepository usando GORM sobre SQLite
type AccesoSQLite struct {
	db *gorm.DB
}

// NewAccesoSQLite construye el repositorio de accesos inyectándole la conexión ya abierta
func NewAccesoSQLite(db *gorm.DB) *AccesoSQLite {
	return &AccesoSQLite{db: db}
}

// ACCESOS

func (s *AccesoSQLite) ListarAccesos() []models.AccesoVehiculo {
	var accesos []models.AccesoVehiculo
	s.db.Find(&accesos)
	return accesos
}

func (s *AccesoSQLite) BuscarAccesoPorID(id uint) (models.AccesoVehiculo, bool) {
	var a models.AccesoVehiculo
	if err := s.db.First(&a, id).Error; err != nil {
		return models.AccesoVehiculo{}, false
	}
	return a, true
}

func (s *AccesoSQLite) CrearAcceso(a models.AccesoVehiculo) models.AccesoVehiculo {
	s.db.Create(&a)
	return a
}

func (s *AccesoSQLite) BorrarAcceso(id uint) bool {
	res := s.db.Delete(&models.AccesoVehiculo{}, id)
	return res.RowsAffected > 0
}
