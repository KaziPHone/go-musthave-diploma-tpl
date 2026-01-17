package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/KaziPHone/go-musthave-diploma-tpl/pkg/helpers"
	"github.com/KaziPHone/go-musthave-diploma-tpl/pkg/user"
)

func (h *Handler) WithdrawHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(user.UserIDKey).(int)

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

	fmt.Println(userBalance.Current, err)

	w.Header().Set("Content-Type", "application/json")
}
