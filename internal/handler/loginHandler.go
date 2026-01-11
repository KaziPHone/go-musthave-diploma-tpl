package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/KaziPHone/go-musthave-diploma-tpl/pkg/crypto"
	"github.com/KaziPHone/go-musthave-diploma-tpl/pkg/user"
	"github.com/golang-jwt/jwt/v4"
)

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {

	var req user.LoginRequest

	isHashed := r.Header.Get("X-Password-Format") == "sha256"

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if req.Login == "" || req.Password == "" {
		http.Error(w, "Login and password are required", http.StatusBadRequest)
		return
	}

	userID, err := h.loginUser(req, isHashed)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid login/password", http.StatusUnauthorized)
		return
	}

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &user.Claims{
		UserID: userID,
		Login:  req.Login,
		RegisteredClaims: jwt.RegisteredClaims{ // ← изменилось имя поля
			ExpiresAt: jwt.NewNumericDate(expirationTime), // ← новый способ
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.JwtKey)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
		Path:    "/",
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}

func (h *Handler) loginUser(req user.LoginRequest, isHashed bool) (int, error) {
	query := `SELECT id FROM users WHERE login = $1 AND password = $2`
	pass := req.Password
	if !isHashed {
		pass = crypto.HashString(req.Password)
	}
	return h.Storage.DbStorage.CountRows(query, req.Login, pass)
}
