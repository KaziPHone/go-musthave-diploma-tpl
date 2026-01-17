package handler

import (
	"database/sql"
	"encoding/json"

	"net/http"

	"github.com/KaziPHone/go-musthave-diploma-tpl/pkg/user"
)

func (h *Handler) GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(user.UserIDKey).(int)

	query := `SELECT number, status, accrual, uploaded_at FROM orders WHERE user_id = $1 ORDER BY uploaded_at DESC`
	rows, err := h.Storage.DBStorage.GetRows(r.Context(), query, userID)

	w.Header().Set("Content-Type", "application/json")

	if rows.Err() != nil || err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var orders []user.Order
	hasOrders := false

	for rows.Next() {
		hasOrders = true
		var order user.Order
		var accrual sql.NullFloat64
		err := rows.Scan(&order.Number, &order.Status, &accrual, &order.UploadedAt)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if accrual.Valid {
			order.Accrual = &accrual.Float64
		}
		orders = append(orders, order)
	}

	if !hasOrders {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}
