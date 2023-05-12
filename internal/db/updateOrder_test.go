package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDB_UpdateOrder(t *testing.T) {
	orders := [3]string{"111111111111111", "2222222222222222", "33333333333333333"}
	accrual := float64(100)
	ctx := context.Background()
	db, err := NewDB(ctx, "user=postgres sslmode=disable")
	require.NoError(t, err)
	user, err := db.Register(ctx, User{Login: "TestDB_UpdateOrderLogin", Password: "TestDB_UpdateOrderPassword"})
	require.NoError(t, err)
	_, err = db.db.ExecContext(ctx, insertOrderSQL, user, orders[0], `NEW`, `now`)
	require.NoError(t, err)
	_, err = db.db.ExecContext(ctx, insertOrderSQL, user, orders[1], `NEW`, `now`)
	require.NoError(t, err)
	type args struct {
		ctx    context.Context
		userID int64
		order  *Order
	}
	tests := []struct {
		name    string
		db      *DB
		args    args
		wantErr bool
	}{
		{
			name: "OK: Simple status update",
			db:   db,
			args: args{
				ctx:    ctx,
				userID: user,
				order:  &Order{Number: orders[0], Status: "REGISTERED"},
			},
			wantErr: false,
		},
		{
			name: "OK: Accrual update (insert new record in balance)",
			db:   db,
			args: args{
				ctx:    ctx,
				userID: user,
				order:  &Order{Number: orders[0], Status: "PROCESSED", Accrual: &accrual},
			},
			wantErr: false,
		},
		{
			name: "OK: Accrual update (update balance)",
			db:   db,
			args: args{
				ctx:    ctx,
				userID: user,
				order:  &Order{Number: orders[1], Status: "PROCESSED", Accrual: &accrual},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.db.UpdateOrder(tt.args.ctx, tt.args.userID, tt.args.order); (err != nil) != tt.wantErr {
				t.Errorf("DB.UpdateOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	_, err = db.db.ExecContext(ctx, "DELETE FROM orders WHERE user_id = $1", user)
	require.NoError(t, err)
	_, err = db.db.ExecContext(ctx, "DELETE FROM balances WHERE user_id = $1", user)
	require.NoError(t, err)
	_, err = db.db.ExecContext(ctx, "DELETE FROM users WHERE user_id = $1", user)
	require.NoError(t, err)
}
