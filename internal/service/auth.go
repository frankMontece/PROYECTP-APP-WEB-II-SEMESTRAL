package service

import (
	"condominio-api/internal/models"
	"condominio-api/internal/storage"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secretJWT = []byte("cambiar_en_produccion")

const duracionToken = 24 * time.Hour

// Claims estructura del JWT
type Claims struct {
	UsuarioID int `json:"usuario_id"`
	jwt.RegisteredClaims
}

// AuthService maneja autenticación y autorización
type AuthService struct {
	repo storage.UserRepository
}

// NewAuthService crea una nueva instancia del servicio de autenticación
func NewAuthService(repo storage.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

// Registrar crea un nuevo usuario
func (s *AuthService) Registrar(email, password string) (models.Usuario, error) {
	email = strings.TrimSpace(email)
	if email == "" || password == "" {
		return models.Usuario{}, ErrCredencialesInvalidas
	}

	// Verificar si el email ya existe
	if _, existe := s.repo.BuscarUsuarioPorEmail(email); existe {
		return models.Usuario{}, ErrEmailEnUso
	}

	// Hashear contraseña
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.Usuario{}, err
	}

	// Crear usuario
	usuario, err := s.repo.CrearUsuario(models.Usuario{
		Email:        email,
		PasswordHash: string(hash),
	})
	if err != nil {
		return models.Usuario{}, err
	}

	return usuario, nil
}

// Login autentica un usuario y genera un JWT
func (s *AuthService) Login(email, password string) (string, error) {
	u, existe := s.repo.BuscarUsuarioPorEmail(strings.TrimSpace(email))
	if !existe {
		return "", ErrCredencialesInvalidas
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", ErrCredencialesInvalidas
	}

	return s.generarToken(u)
}

// generarToken crea un JWT para el usuario
func (s *AuthService) generarToken(u models.Usuario) (string, error) {
	claims := &Claims{
		UsuarioID: u.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duracionToken)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretJWT)
}

// ValidarToken verifica un JWT y retorna el ID del usuario
func (s *AuthService) ValidarToken(token string) (int, error) {
	parsed, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrCredencialesInvalidas
		}
		return secretJWT, nil
	})

	if err != nil || !parsed.Valid {
		return 0, ErrCredencialesInvalidas
	}

	claims, ok := parsed.Claims.(*Claims)
	if !ok {
		return 0, ErrCredencialesInvalidas
	}

	return claims.UsuarioID, nil
}
