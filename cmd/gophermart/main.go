package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/NevostruevK/GopherMart.git/internal/db"
	"github.com/NevostruevK/GopherMart.git/internal/util/logger"
)

func main() {
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	lg := logger.NewLogger("main : ", log.LstdFlags|log.Lshortfile)

	db, err := db.NewDB(ctx, "user=postgres sslmode=disable")
	if err != nil {
		lg.Println(err)
		return
	}
	defer func() {
		err := db.Close()
		lg.Println(err)
	}()
	<-gracefulShutdown
}
