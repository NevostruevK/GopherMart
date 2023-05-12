package db

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const errTooLong = `pq: value too long`

func TestDB_Register(t *testing.T) {
	newId := make([]int64, 0, 1)
	ctx := context.Background()
	db, err := NewDB(ctx, "user=postgres sslmode=disable")
	require.NoError(t, err)
	type args struct {
		ctx context.Context
		u   User
	}
	tests := []struct {
		name      string
		db        *DB
		args      args
		idNotZero bool
		wantErr   bool
		waitErr   string
	}{
		{
			name: "add normal user",
			db:   db,
			args: args{
				ctx: ctx,
				u:   User{Login: "testUser123", Password: "testUser123Password"}},
			idNotZero: true,
			wantErr:   false,
			waitErr:   "",
		},
		{
			name: "add same login ",
			db:   db,
			args: args{
				ctx: ctx,
				u:   User{Login: "testUser123", Password: "testUser123AnotherPassword"}},
			idNotZero: false,
			wantErr:   true,
			waitErr:   errDuplicateLogin,
		},
		{
			name: "add big login ",
			db:   db,
			args: args{
				ctx: ctx,
				u:   User{Login: "testUser123veryveryverybig logginnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn", Password: "testUser123AnotherPassword"}},
			idNotZero: false,
			wantErr:   true,
			waitErr:   errTooLong,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.db.Register(tt.args.ctx, tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				if !strings.Contains(err.Error(), tt.waitErr) {
					t.Errorf("DB.Register() error = %v, waitErr %v", err, tt.waitErr)
					return
				}
			}
			if got > 0 != tt.idNotZero {
				t.Errorf("DB.Register() = %v, want %v", got, tt.idNotZero)
			}
			if got > 0 {
				newId = append(newId, got)
			}
		})
	}
	for _, id := range newId {
		_, err = db.db.ExecContext(ctx, "DELETE FROM users WHERE user_id = $1", id)
		require.NoError(t, err)
	}

}
