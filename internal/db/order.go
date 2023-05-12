package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

const errDuplicateOrder = `pq: duplicate key value violates unique constraint "orders_number_key"`
const initialOrdersCount = 8

type Order struct {
	Number   string    `json:"number"`                // Номер заказа
	Status   string    `json:"status"`                // Статус заказа
	Accrual  *float64  `json:"accrual,omitempty"`     // Начислено баллов
	Uploaded time.Time `json:"uploaded_at,omitempty"` // Время загрузки заказа
}

func (o Order) String() string {
	if o.Accrual == nil {
		return fmt.Sprintf("Number:%s Status:%10s Uploaded:%v", o.Number, o.Status, o.Uploaded)
	}
	return fmt.Sprintf("Number:%s Status:%10s Accrual:%f Uploaded:%v", o.Number, o.Status, *o.Accrual, o.Uploaded)
}

func (db *DB) PostOrder(ctx context.Context, userID int64, order string) (int64, error) {
	if _, err := db.db.ExecContext(ctx, insertOrderSQL, userID, order, `NEW`, `now`); err != nil {
		db.lg.Println(err)
		if !strings.Contains(err.Error(), errDuplicateOrder) {
			return 0, err
		}
		var id int64
		if err := db.db.QueryRowContext(ctx, getUserForOrderSQL, order).Scan(&id); err != nil {
			db.lg.Printf("ERROR : getUserForOrder %s %v\n", order, err)
			return 0, err
		}
		return id, fmt.Errorf(errDuplicateOrder)
	}
	return userID, nil
}

func (db *DB) GetOrders(ctx context.Context, userID int64) ([]Order, error) {
	orders := make([]Order, 0, initialOrdersCount)
	rows, err := db.db.QueryContext(ctx, getOrdersSQL, userID)
	if err != nil {
		db.lg.Println(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		order := Order{}
		accrual := sql.NullFloat64{}
		if err = rows.Scan(&order.Number, &order.Status, &accrual, &order.Uploaded); err != nil {
			db.lg.Println(err)
			continue
		}
		if accrual.Valid {
			order.Accrual = &accrual.Float64
		}
		orders = append(orders, order)
	}
	return orders, rows.Err()
}
