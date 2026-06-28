package storage

import (
	"condominio-api/internal/models"

	"gorm.io/gorm"
)

// AreaSQLite implementa AreaSocialRepository usando GORM sobre SQLite
type AreaSQLite struct {
	db *gorm.DB
}

// NewAreaSQLite construye el repositorio de áreas inyectándole la conexión ya abierta
func NewAreaSQLite(db *gorm.DB) *AreaSQLite {
	return &AreaSQLite{db: db}
}

// Verificación en compilación: si AreaSQLite no implementa AreaSocialRepository
var _ AreaSocialRepository = (*AreaSQLite)(nil)

// ─── ÁREAS SOCIALES ──────────────────────────────────────────────────────────

func (s *AreaSQLite) ListarAreas() []models.AreaSocial {
	var areas []models.AreaSocial
	s.db.Find(&areas)
	return areas
}

func (s *AreaSQLite) BuscarAreaPorID(id uint) (models.AreaSocial, bool) {
	var a models.AreaSocial
	if err := s.db.First(&a, id).Error; err != nil {
		return models.AreaSocial{}, false
	}
	return a, true
}

func (s *AreaSQLite) CrearArea(a models.AreaSocial) models.AreaSocial {
	s.db.Create(&a)
	return a
}

func (s *AreaSQLite) ActualizarArea(id uint, datos models.AreaSocial) (models.AreaSocial, bool) {
	var existente models.AreaSocial
	if err := s.db.First(&existente, id).Error; err != nil {
		return models.AreaSocial{}, false
	}
	datos.ID = id
	s.db.Save(&datos)
	return datos, true
}

func (s *AreaSQLite) BorrarArea(id uint) bool {
	res := s.db.Delete(&models.AreaSocial{}, id)
	return res.RowsAffected > 0
}
