package db

import (
	"context"
)

const ErrDuplicateLogin = `pq: duplicate key value violates unique constraint "users_login_key"`

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (db *DB) Register(ctx context.Context, u User) (uint64, error) {
	if _, err := db.db.ExecContext(ctx, insertUserSQL, u.Login, u.Password); err != nil {
		db.lg.Println(err)
		return 0, err
	}
	return db.Login(ctx, u)
}

func (db *DB) Login(ctx context.Context, u User) (uint64, error) {
	var id uint64
	if err := db.db.QueryRowContext(ctx, getUserSQL, u.Login, u.Password).Scan(&id); err != nil {
		db.lg.Printf("ERROR : getUser %s %v\n", u.Login, err)
		return 0, err
	}
	return id, nil
}
