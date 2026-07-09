package storage

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"condominio-api/internal/models"
)

// SembrarUsuarios inserta usuarios de ejemplo solo si la tabla está vacía.
// Debe llamarse ANTES que cualquier otro Sembrar* que use ResidenteID como
// foreign key (Obligaciones, Multas, Parqueo, Área Social), porque Postgres
// valida esa constraint estrictamente — a diferencia de SQLite, que por
// defecto no la aplica.
func SembrarUsuarios(db *gorm.DB) {
	var n int64
	db.Model(&models.Usuario{}).Count(&n)
	if n > 0 {
		return
	}

	usuarios := []models.Usuario{
		{Email: "residente1@test.com", PasswordHash: "seed", CreadoEn: time.Now()},
		{Email: "residente2@test.com", PasswordHash: "seed", CreadoEn: time.Now()},
		{Email: "residente3@test.com", PasswordHash: "seed", CreadoEn: time.Now()},
	}
	db.Create(&usuarios)

	fmt.Printf("✅ Datos de ejemplo sembrados en Usuarios: %d usuarios\n", len(usuarios))
}
