package handler

import (
	"encoding/json"
	"net/http"

	"github.com/KaziPHone/go-musthave-diploma-tpl/pkg/helpers"
	"github.com/KaziPHone/go-musthave-diploma-tpl/pkg/user"
)

func (h *Handler) GetBalanceHandler(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(user.UserIDKey).(int)

	w.Header().Set("Content-Type", "application/json")

	var userBalance user.Balance
	query := `SELECT user_id, current, withdrawn FROM user_balance WHERE user_id = $1`
	row := h.Storage.DBStorage.InsertWithReturning(query, userID)

	userBalance, err := helpers.GetUserBalance(row)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userBalance)
}
