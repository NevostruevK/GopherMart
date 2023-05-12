package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/NevostruevK/GopherMart.git/internal/util/logger"
	_ "github.com/lib/pq"
)

type DB struct {
	db   *sql.DB
	lg   *log.Logger
	init bool
}

func NewDB(ctx context.Context, connStr string) (*DB, error) {
	db := &DB{db: nil, lg: logger.NewLogger("postgres : ", log.LstdFlags|log.Lshortfile), init: false}
	if connStr == "" {
		msg := "DATABASE_URI is empty, database wasn't initialized"
		db.lg.Println(msg)
		return db, fmt.Errorf(msg)
	}
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		db.lg.Println(err)
		return db, err
	}
	db.db = conn

	if err = db.initTables(ctx); err != nil {
		db.lg.Println(err)
		return db, err
	}
	db.init = true
	return db, nil
}

func (db *DB) initTables(ctx context.Context) error {
	if _, err := db.db.ExecContext(ctx, createUsersSQL); err != nil {
		return err
	}
	if _, err := db.db.ExecContext(ctx, createBalancesSQL); err != nil {
		return err
	}
	if _, err := db.db.ExecContext(ctx, createOrdersSQL); err != nil {
		return err
	}
	if _, err := db.db.ExecContext(ctx, createWithdrawalSQL); err != nil {
		return err
	}
	return nil
}

func (db *DB) Close() error {
	if !db.init {
		return fmt.Errorf(" Can't close DB : DataBase wasn't initiated")
	}
	if err := db.db.Close(); err != nil {
		return fmt.Errorf(" Can't close DB %w", err)
	}
	return nil
}
