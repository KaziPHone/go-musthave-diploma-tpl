package helpers

import (
	"database/sql"

	"github.com/KaziPHone/go-musthave-diploma-tpl/pkg/user"
)

func GetUserBalance(row *sql.Row) (user.Balance, error) {
	var userBalance user.Balance
	if row.Err() != nil {
		return user.Balance{}, row.Err()
	}
	if err := row.Scan(&userBalance.UserId, &userBalance.Current, &userBalance.Withdrawn); err != nil {
		return user.Balance{}, err
	}
	return userBalance, nil
}
