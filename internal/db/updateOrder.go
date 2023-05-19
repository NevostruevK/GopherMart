package db

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/NevostruevK/GopherMart.git/internal/client/task"
)

func (db *DB) UpdateOrder(ctx context.Context, userID uint64, order *task.Order) error {
	if order.Accrual == nil {
		return db.UpdateStatusOrder(ctx, userID, order)
	}
	return db.UpdateAccrual(ctx, userID, order)
}

func (db *DB) UpdateStatusOrder(ctx context.Context, userID uint64, order *task.Order) error {
	if _, err := db.db.ExecContext(ctx, updateOrderStatusSQL, order.Number, order.Status); err != nil {
		db.lg.Println(err)
		return err
	}
	return nil
}

func (db *DB) UpdateAccrual(ctx context.Context, userID uint64, order *task.Order) error {
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
	if _, err := tx.ExecContext(ctx, updateOrderAccrualSQL, order.Number, order.Status, *order.Accrual); err != nil {
		db.lg.Println(err)
		return err
	}
	var current float64
	if err := tx.QueryRowContext(ctx, getCurrentBalanceSQL, userID).Scan(&current); err != nil {
		if !strings.Contains(err.Error(), errNoBalance) {
			db.lg.Println(err)
			return err
		}
		if _, err = tx.ExecContext(ctx, insertBalanceSQL, userID, *order.Accrual, 0); err != nil {
			db.lg.Println(err)
			return err
		}
		return tx.Commit()
	}
	if _, err = tx.ExecContext(ctx, updateCurrentBalanceSQL, userID, current+*order.Accrual); err != nil {
		db.lg.Println(err)
		return err
	}
	return tx.Commit()
}
