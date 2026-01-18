package handler

import (
	"encoding/json"
	"net/http"

	"github.com/KaziPHone/go-musthave-diploma-tpl/pkg/user"
)

func (h *Handler) GetWithdrawalsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(user.UserIDKey).(int)

	w.Header().Set("Content-Type", "application/json")

	query := `SELECT number, status, accrual, uploaded_at FROM balance_operationse WHERE user_id = $1 ORDER BY uploaded_at DESC`
	rows, err := h.Storage.DBStorage.GetRows(r.Context(), query, userID)

	if err != nil || rows == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var operations []user.Operation
	for rows.Next() {
		var operation user.Operation
		err := rows.Scan(&operation.ID, &operation.UserID, &operation.Type, &operation.Amount, &operation.Order, &operation.ProcessedAt)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		operations = append(operations, operation)
	}

	if len(operations) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(operations)
}
