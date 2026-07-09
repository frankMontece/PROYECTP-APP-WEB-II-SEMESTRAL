package storage

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"condominio-api/internal/models"
)

// SembrarUsuarios inserta usuarios de ejemplo solo si la tabla está vacía.
// Debe llamarse ANTES que cualquier otro Sembrar* que use ResidenteID como
// foreign key (Obligaciones, Multas, Parqueo, Área Social), porque Postgres
// valida esa constraint estrictamente — a diferencia de SQLite, que por
// defecto no la aplica.
//
// residente1@test.com queda como "admin" para tener un usuario con permisos
// completos listo en la demo, sin pasos manuales.
func SembrarUsuarios(db *gorm.DB) {
	var n int64
	db.Model(&models.Usuario{}).Count(&n)
	if n > 0 {
		return
	}

	// Password real para poder hacer login en la demo: "password123" para todos.
	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("⚠️  error generando hash de seed: %v\n", err)
		return
	}

	usuarios := []models.Usuario{
		{Email: "residente1@test.com", PasswordHash: string(hash), Rol: "admin", CreadoEn: time.Now()},
		{Email: "residente2@test.com", PasswordHash: string(hash), Rol: "residente", CreadoEn: time.Now()},
		{Email: "residente3@test.com", PasswordHash: string(hash), Rol: "residente", CreadoEn: time.Now()},
	}
	db.Create(&usuarios)

	fmt.Printf("✅ Datos de ejemplo sembrados en Usuarios: %d usuarios (residente1@test.com = admin)\n", len(usuarios))
}
