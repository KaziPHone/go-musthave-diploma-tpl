package handler

import (
	"encoding/json"
	"net/http"

	"github.com/KaziPHone/go-musthave-diploma-tpl/pkg/crypto"
	"github.com/KaziPHone/go-musthave-diploma-tpl/pkg/user"
)

func (h *Handler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {

	isHashed := r.Header.Get("X-Password-Format") == "sha256"

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}

	var req user.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Login == "" || req.Password == "" {
		http.Error(w, "Login and password are required", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	userIsRegistred, err := h.userIsRegistred(req.Login)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if userIsRegistred {
		w.WriteHeader(http.StatusConflict)
		return
	}

	err = h.registrationUser(req, isHashed)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully",
	})

}

func (h *Handler) registrationUser(req user.RegisterRequest, isHashed bool) error {
	query := `INSERT INTO users (login, password) VALUES ($1, $2)`
	pass := req.Password
	if !isHashed {
		pass = crypto.HashString(req.Password) // хеширование пароля
	}
	err := h.Storage.DbStorage.Insert(query, req.Login, pass)
	return err
}

func (h *Handler) userIsRegistred(login string) (bool, error) {
	query := `SELECT COUNT(login) FROM users WHERE login = $1`
	result, err := h.Storage.DbStorage.CountRows(query, login)
	return result > 0, err
}
