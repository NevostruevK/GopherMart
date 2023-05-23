package db

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const errDuplicateWithdrawal = `pq: duplicate key value violates unique constraint "withdrawal_number_key"`

func TestDB_PostWithdrawnOrder(t *testing.T) {
	order := "123456789123214214"
	users := make([]uint64, 0, 2)
	withdrawal := float64(100)
	bigWithdrawal := float64(1000000)
	ctx := context.Background()
	db, err := NewDB(ctx, "user=postgres sslmode=disable")
	require.NoError(t, err)
	userID, err := db.Register(ctx, User{Login: "testUser12345678901", Password: "testUser12345678901Password"})
	require.NoError(t, err)
	_, err = db.db.ExecContext(ctx, insertBalanceSQL, userID, 1000, 10000)
	require.NoError(t, err)
	users = append(users, userID)
	userID, err = db.Register(ctx, User{Login: "testUser1234567890", Password: "testUser1234567890Password"})
	require.NoError(t, err)
	users = append(users, userID)

	type args struct {
		ctx    context.Context
		userID uint64
		order  *WithdrawnOrder
	}
	tests := []struct {
		name    string
		db      *DB
		args    args
		wantErr bool
		waitErr string
	}{
		{
			name: "ERROR: not enough money (no record in balance)",
			db:   db,
			args: args{
				ctx:    ctx,
				userID: users[1],
				order:  &WithdrawnOrder{Number: order, Withdrawn: withdrawal},
			},
			wantErr: true,
			waitErr: ErrNotEnoughFounds,
		},
		{
			name: "ERROR: not enough money (there is a record in balance)",
			db:   db,
			args: args{
				ctx:    ctx,
				userID: users[0],
				order:  &WithdrawnOrder{Number: order, Withdrawn: bigWithdrawal},
			},
			wantErr: true,
			waitErr: ErrNotEnoughFounds,
		},
		{
			name: "OK: add normal withdrawal",
			db:   db,
			args: args{
				ctx:    ctx,
				userID: users[0],
				order:  &WithdrawnOrder{Number: order, Withdrawn: withdrawal},
			},
			wantErr: false,
			waitErr: "",
		},
		{
			name: "ERROR: duplicate order",
			db:   db,
			args: args{
				ctx:    ctx,
				userID: users[0],
				order:  &WithdrawnOrder{Number: order, Withdrawn: withdrawal},
			},
			wantErr: true,
			waitErr: errDuplicateWithdrawal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.db.PostWithdrawal(tt.args.ctx, tt.args.userID, tt.args.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.PostWithdrawal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				if !strings.Contains(err.Error(), tt.waitErr) {
					t.Errorf("DB.PostWithdrawal() error = %v, waitErr %v", err, tt.waitErr)
				}
			}

		})
	}
	_, err = db.db.ExecContext(ctx, "DELETE FROM withdrawal WHERE user_id = $1", users[0])
	require.NoError(t, err)
	_, err = db.db.ExecContext(ctx, "DELETE FROM balances WHERE user_id = $1", users[0])
	require.NoError(t, err)
	for _, id := range users {
		_, err = db.db.ExecContext(ctx, "DELETE FROM users WHERE user_id = $1", id)
		require.NoError(t, err)
	}
}
