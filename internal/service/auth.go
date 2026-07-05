package service

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"condominio-api/internal/models"
	"condominio-api/internal/storage"
)

// Claims define la carga útil del token JWT.
type Claims struct {
	UsuarioID int `json:"usuario_id"`
	jwt.RegisteredClaims
}

// AuthService encapsula toda la lógica de autenticación.
type AuthService struct {
	repo          storage.UserRepository
	secretJWT     []byte
	duracionToken time.Duration
}

// Constructor.
func NewAuthService(
	repo storage.UserRepository,
	secret []byte,
	duracion time.Duration,
) *AuthService {
	return &AuthService{
		repo:          repo,
		secretJWT:     secret,
		duracionToken: duracion,
	}
}

// Registrar valida los datos, hashea la contraseña y crea el usuario.
func (s *AuthService) Registrar(email, password string) (models.Usuario, error) {
	email = strings.TrimSpace(email)
	password = strings.TrimSpace(password)

	if email == "" || password == "" {
		return models.Usuario{}, ErrCredencialesInvalidas
	}
	if _, existe := s.repo.BuscarUsuarioPorEmail(email); existe {
		return models.Usuario{}, ErrEmailEnUso
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.Usuario{}, err
	}

	return s.repo.CrearUsuario(models.Usuario{
		Email:        email,
		PasswordHash: string(hash),
	})
}

// Login verifica las credenciales y devuelve un JWT.
func (s *AuthService) Login(email, password string) (string, error) {
	email = strings.TrimSpace(email)
	password = strings.TrimSpace(password)

	if email == "" || password == "" {
		return "", ErrCredencialesInvalidas
	}

	u, existe := s.repo.BuscarUsuarioPorEmail(email)
	if !existe {
		return "", ErrCredencialesInvalidas
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", ErrCredencialesInvalidas
	}

	return s.GenerarToken(u)
}

// Genera un JWT firmado.
func (s *AuthService) GenerarToken(u models.Usuario) (string, error) {
	claims := &Claims{
		UsuarioID: u.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.duracionToken)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.secretJWT)
}

// Verifica un JWT.
func (s *AuthService) ValidarToken(tokenStr string) (int, error) {
	parsedToken, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrCredencialesInvalidas
		}

		return s.secretJWT, nil
	})

	if err != nil || !parsedToken.Valid {
		return 0, ErrCredencialesInvalidas
	}

	claims, ok := parsedToken.Claims.(*Claims)
	if !ok {
		return 0, ErrCredencialesInvalidas
	}

	return claims.UsuarioID, nil
}
