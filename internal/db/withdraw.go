package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

const ErrNotEnoughFounds = `there are not enough founds`
const initialWithdrawnOrdersCount = 4

type WithdrawnOrder struct {
	Number    string     `json:"order"`                  // Номер заказа
	Withdrawn float64    `json:"sum"`                    // Списано баллов
	Uploaded  *time.Time `json:"processed_at,omitempty"` // Время загрузки заказа
}

func (o WithdrawnOrder) String() string {
	if o.Uploaded == nil {
		return fmt.Sprintf("Number:%s Withdrawn:%f", o.Number, o.Withdrawn)
	}
	return fmt.Sprintf("Number:%s Withdrawn:%f Uploaded:%v", o.Number, o.Withdrawn, *o.Uploaded)
}

func (db *DB) GetWithdrawals(ctx context.Context, userID uint64) ([]WithdrawnOrder, error) {
	orders := make([]WithdrawnOrder, 0, initialWithdrawnOrdersCount)
	rows, err := db.db.QueryContext(ctx, getWithdrawalsSQL, userID)
	if err != nil {
		db.lg.Println(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		order := WithdrawnOrder{}
		if err = rows.Scan(&order.Number, &order.Withdrawn, &order.Uploaded); err != nil {
			db.lg.Println(err)
			continue
		}
		orders = append(orders, order)
	}
	return orders, rows.Err()
}

func (db *DB) PostWithdrawal(ctx context.Context, userID uint64, order *WithdrawnOrder) error {
	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		db.lg.Println(err)
		return err
	}
	defer func() {
		err := tx.Rollback()
		if !errors.Is(err, sql.ErrTxDone) {
			db.lg.Println(err)
		}
	}()
	b := &Balance{}
	if err := tx.QueryRowContext(ctx, getBalanceSQL, userID).Scan(&b.Current, &b.Withdrawn); err != nil {
		if !strings.Contains(err.Error(), errNoBalance) {
			db.lg.Printf("ERROR : getBalance %d %v\n", userID, err)
			return err
		}
		b = NewBalance(0, 0)
	}
	if b.Current < order.Withdrawn {
		return fmt.Errorf(ErrNotEnoughFounds)
	}
	if _, err := tx.ExecContext(ctx, insertWithdrawalSQL, userID, order.Number, order.Withdrawn, `now`); err != nil {
		db.lg.Println(err)
		return err
	}
	if _, err := tx.ExecContext(ctx, updateBalanceSQL, userID, b.Current-order.Withdrawn, b.Withdrawn+order.Withdrawn); err != nil {
		db.lg.Println(err)
		return err
	}
	return tx.Commit()
}
