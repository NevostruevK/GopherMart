package db

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const errForeignKeyOrder = `pq: insert or update on table "orders" violates foreign key constraint "orders_user_id_fkey"`

func TestDB_InsertOrder(t *testing.T) {
	order := "123456789012312321321"
	anotherOrder := "123456789012312321321123124"
	users := make([]int64, 0, 2)
	ctx := context.Background()
	db, err := NewDB(ctx, "user=postgres sslmode=disable")
	require.NoError(t, err)
	userID, err := db.Register(ctx, User{Login: "testUser123456789", Password: "testUser123456789Password"})
	require.NoError(t, err)
	users = append(users, userID)
	userID, err = db.Register(ctx, User{Login: "testUser1234567890", Password: "testUser1234567890Password"})
	require.NoError(t, err)
	users = append(users, userID)

	type args struct {
		ctx    context.Context
		userID int64
		order  string
	}
	tests := []struct {
		name    string
		db      *DB
		args    args
		wantID  int64
		wantErr bool
		waitErr string
	}{
		{
			name: "add normal order",
			db:   db,
			args: args{
				ctx:    ctx,
				userID: users[0],
				order:  order,
			},
			wantID:  users[0],
			wantErr: false,
			waitErr: "",
		},
		{
			name: "add the same order for the same user",
			db:   db,
			args: args{
				ctx:    ctx,
				userID: users[0],
				order:  order,
			},
			wantID:  users[0],
			wantErr: true,
			waitErr: errDuplicateOrder,
		},
		{
			name: "add the same order for another user",
			db:   db,
			args: args{
				ctx:    ctx,
				userID: users[1],
				order:  order,
			},
			wantID:  users[0],
			wantErr: true,
			waitErr: errDuplicateOrder,
		},
		{
			name: "add order for the wrong user",
			db:   db,
			args: args{
				ctx:    ctx,
				userID: userID + 1,
				order:  anotherOrder,
			},
			wantID:  0,
			wantErr: true,
			waitErr: errForeignKeyOrder,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.db.PostOrder(tt.args.ctx, tt.args.userID, tt.args.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.PostOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				if !strings.Contains(err.Error(), tt.waitErr) {
					t.Errorf("DB.PostOrder() error = %v, waitErr %v", err, tt.waitErr)
					return
				}
			}
			if got != tt.wantID {
				t.Errorf("DB.PostOrder() = %v, want %v", got, tt.wantID)
			}
		})
	}
	_, err = db.db.ExecContext(ctx, "DELETE FROM orders WHERE user_id = $1", users[0])
	require.NoError(t, err)
	for _, id := range users {
		_, err = db.db.ExecContext(ctx, "DELETE FROM users WHERE user_id = $1", id)
		require.NoError(t, err)
	}
}
