package db

import (
	"context"
	"fmt"
	"strings"
)

const errNoBalance = `sql: no rows in result set`

type Balance struct {
	Current   float64 `json:"current"`   // Баланс баллов
	Withdrawn float64 `json:"withdrawn"` // Списано баллов
}

func (b Balance) String() string{
	return fmt.Sprintf("balance %f :%f ",b.Current, b.Withdrawn)
}

func NewBalance(current, withdrawn float64) *Balance {
	return &Balance{current, withdrawn}
}

func (db *DB) GetBalance(ctx context.Context, userID uint64) (*Balance, error) {
	b := Balance{}
	if err := db.db.QueryRowContext(ctx, getBalanceSQL, userID).Scan(&b.Current, &b.Withdrawn); err != nil {
		if !strings.Contains(err.Error(), errNoBalance) {
			db.lg.Printf("ERROR : getBalance %d %v\n", userID, err)
			return nil, err
		}
		return NewBalance(0, 0), nil
	}
	return &b, nil
}
