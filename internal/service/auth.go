package service

import (
	"condominio-api/internal/models"
	"condominio-api/internal/storage"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Valores por defecto para desarrollo. En produccion se inyectan via Options
// desde la configuracion cargada del .env (ver internal/config).
const (
	secretoPorDefecto  = "condominio-secreto-solo-dev"
	duracionPorDefecto = 24 * time.Hour
)

type Claims struct {
	UsuarioID int `json:"usuario_id"`
	jwt.RegisteredClaims
}

// AuthService maneja autenticación y autorización.
//
// Antes: secreto y duracion eran una VARIABLE y una CONSTANTE globales.
// Ahora son campos del struct, configurables por el patron Options.
type AuthService struct {
	repo     storage.UserRepository
	secreto  []byte
	duracion time.Duration
}

// AuthOption configura un AuthService en su construccion.
type AuthOption func(*AuthService)

// WithSecreto inyecta la clave de firma del JWT. Si recibe vacio, conserva el default.
func WithSecreto(secreto []byte) AuthOption {
	return func(a *AuthService) {
		if len(secreto) > 0 {
			a.secreto = secreto
		}
	}
}

// WithDuracionToken inyecta la validez del token. Si recibe <= 0, conserva el default.
func WithDuracionToken(d time.Duration) AuthOption {
	return func(a *AuthService) {
		if d > 0 {
			a.duracion = d
		}
	}
}

// NewAuthService crea una nueva instancia del servicio de autenticación.
// opts es variadico: NewAuthService(repo) sigue compilando sin cambios.
func NewAuthService(repo storage.UserRepository, opts ...AuthOption) *AuthService {
	s := &AuthService{
		repo:     repo,
		secreto:  []byte(secretoPorDefecto),
		duracion: duracionPorDefecto,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Registrar crea un nuevo usuario (sin cambios de lógica).
func (s *AuthService) Registrar(email, password string) (models.Usuario, error) {
	email = strings.TrimSpace(email)
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

// Login autentica un usuario y genera un JWT (sin cambios de lógica).
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

// generarToken usa s.duracion y s.secreto en vez de las variables globales.
func (s *AuthService) generarToken(u models.Usuario) (string, error) {
	claims := &Claims{
		UsuarioID: u.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.duracion)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secreto)
}

// ValidarToken usa s.secreto en vez de la variable global.
func (s *AuthService) ValidarToken(token string) (int, error) {
	parsed, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrCredencialesInvalidas
		}
		return s.secreto, nil
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
