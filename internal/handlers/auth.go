package handlers

import (
	"encoding/json"
	"net/http"
)

// credenciales es el body que reciben Register y Login.
// Se define aquí para no exponer structs internos del modelo.
type credenciales struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Registrar atiende POST /api/v1/auth/register.
func (s *Server) Registrar(w http.ResponseWriter, r *http.Request) {
	var creds credenciales
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}

	usuario, err := s.Auth.Registrar(creds.Email, creds.Password)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusCreated, usuario)
}

// Login atiende POST /api/v1/auth/login.
func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	var creds credenciales
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}

	token, err := s.Auth.Login(creds.Email, creds.Password)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"token": token})
}
