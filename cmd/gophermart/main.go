package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/NevostruevK/GopherMart.git/internal/client"
	"github.com/NevostruevK/GopherMart.git/internal/db"
	"github.com/NevostruevK/GopherMart.git/internal/server"
	"github.com/NevostruevK/GopherMart.git/internal/util/logger"
	"github.com/NevostruevK/GopherMart.git/internal/util/option"
)

func main() {
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	ctx, cancel := context.WithCancel(context.Background())
	lg := logger.NewLogger("main : ", log.LstdFlags|log.Lshortfile)
	opt, err := option.GetOptions()
	if err != nil {
		lg.Printf("GetOptions() error %v", err)
	}
	db, err := db.NewDB(ctx, opt.DatabaseURI)
	if err != nil {
		lg.Println(err)
		return
	}
	defer func() {
		err := db.Close()
		lg.Println(err)
	}()
	m := client.NewManager()
	go m.Start(ctx, db, opt.AccrualSystemAddress, 1)
	s, err := server.NewServer(db, opt.RunAddress, m)
	if err != nil {
		lg.Println(err)
		return
	}
	lg.Printf("Start server")
	go func() {
		go lg.Println(s.ListenAndServe())
	}()
	<-gracefulShutdown
	if err = s.Shutdown(ctx); err != nil {
		lg.Printf("ERROR : Server Shutdown error %v", err)
	} else {
		lg.Printf("Server Shutdown ")
	}
	cancel()
}
