package handler

import (
	"encoding/json"
	"net/http"

	"github.com/KaziPHone/go-musthave-diploma-tpl/pkg/user"
)

func (h *Handler) GetBalanceHandler(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(user.UserIDKey).(int)

	w.Header().Set("Content-Type", "application/json")

	var userBalance user.Balance
	query := `SELECT user_id, current, withdrawn FROM user_balance WHERE user_id = $1`
	row := h.Storage.DBStorage.InsertWithReturning(query, userID)

	if row.Err() != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := row.Scan(&userBalance.UserId, &userBalance.Current, &userBalance.Withdrawn); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userBalance)
}
