package storage

import (
	"condominio-api/internal/models"

	"gorm.io/gorm"
)

// UsuarioGORM implementa UserRepository usando GORM
type UsuarioGORM struct {
	db *gorm.DB
}

// NewUsuarioGORM crea una nueva instancia del repositorio de usuarios
func NewUsuarioGORM(db *gorm.DB) *UsuarioGORM {
	return &UsuarioGORM{db: db}
}

// CrearUsuario implementa UserRepository
func (r *UsuarioGORM) CrearUsuario(u models.Usuario) (models.Usuario, error) {
	if err := r.db.Create(&u).Error; err != nil {
		return models.Usuario{}, err
	}
	return u, nil
}

// BuscarUsuarioPorEmail implementa UserRepository
func (r *UsuarioGORM) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	var u models.Usuario
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return models.Usuario{}, false
	}
	return u, true
}

// Verificación en tiempo de compilación
var _ UserRepository = (*UsuarioGORM)(nil)
