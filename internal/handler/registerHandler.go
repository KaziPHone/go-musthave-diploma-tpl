package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/KaziPHone/go-musthave-diploma-tpl/pkg/crypto"
	"github.com/KaziPHone/go-musthave-diploma-tpl/pkg/helpers"
	"github.com/KaziPHone/go-musthave-diploma-tpl/pkg/user"
)

func (h *Handler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {

	isHashed := r.Header.Get("X-Password-Format") == "sha256"

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}

	var req user.UserRequest
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

	userID, err := h.loginUser(req, isHashed)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := h.addUserBalance(userID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie, err := helpers.GetAuthCookie(userID, req.Login, h.JwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully",
	})

}

func (h *Handler) registrationUser(req user.UserRequest, isHashed bool) error {

	pass := req.Password
	if !isHashed {
		pass = crypto.HashString(req.Password) // хеширование пароля
	}

	query := `INSERT INTO users (login, password) VALUES ($1, $2)`
	if err := h.Storage.DBStorage.Insert(query, req.Login, pass); err != nil {
		return err
	}

	return nil
}

func (h *Handler) userIsRegistred(login string) (bool, error) {
	query := `SELECT COUNT(login) FROM users WHERE login = $1`
	result, err := h.Storage.DBStorage.CountRows(query, login)
	return result > 0, err
}

func (h *Handler) addUserBalance(userId int) error {
	query := `INSERT INTO user_balance (user_id, current, withdrawn, updated_at, uploaded_at) 
	VALUES ($1, $2, $3, $4, $5)`
	t := time.Now()
	return h.Storage.DBStorage.Insert(query, userId, 0, 0, t, t)
}
