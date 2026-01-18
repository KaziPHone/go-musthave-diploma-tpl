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

	balance := roundTo4Decimals(userBalance.Current - req.Sum)
	if userBalance.Current <= 0 || balance < 0 {
		http.Error(w, "Internal server error", http.StatusPaymentRequired)
		return
	}

	query = `
		UPDATE orders 
		SET accrual = accrual - $2, 
			processed_at = $3 
		WHERE number = $1
	`
	err = h.Storage.DBStorage.Insert(query, req.Order, balance, time.Now())
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	fmt.Println(balance, err, req.Sum)
	w.WriteHeader(http.StatusOK)

}

func roundTo4Decimals(value float64) float64 {
	return float64(int64(value*10000+0.5)) / 10000
}
