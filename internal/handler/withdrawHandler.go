package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/KaziPHone/go-musthave-diploma-tpl/pkg/helpers"
	"github.com/KaziPHone/go-musthave-diploma-tpl/pkg/user"
)

func (h *Handler) WithdrawHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(user.UserIDKey).(int)

	w.Header().Set("Content-Type", "application/json")

	type request struct {
		Order string  `json:"order"`
		Sum   float64 `json:"sum"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	query := `SELECT user_id, current, withdrawn FROM user_balance WHERE user_id = $1`
	row := h.Storage.DBStorage.InsertWithReturning(query, userID)
	userBalance, err := helpers.GetUserBalance(row)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if userBalance.Current <= 0 {
		http.Error(w, "Internal server error", http.StatusPaymentRequired)
		return
	}

	balance := roundTo4Decimals(userBalance.Current - req.Sum)
	query = `
		UPDATE user_balance 
		SET current = current - $2, 
			updated_at = $3 
		WHERE user_id = $1
	`
	err = h.Storage.DBStorage.Insert(query, userID, balance, time.Now())
	if err != nil {
		http.Error(w, "Internal server error", http.StatusPaymentRequired)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Println(balance, err)

}

func roundTo4Decimals(value float64) float64 {
	return float64(int64(value*10000+0.5)) / 10000
}
