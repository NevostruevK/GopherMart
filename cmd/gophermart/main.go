package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/NevostruevK/GopherMart.git/internal/db"
	"github.com/NevostruevK/GopherMart.git/internal/server"
	"github.com/NevostruevK/GopherMart.git/internal/util/logger"
	"github.com/NevostruevK/GopherMart.git/internal/util/option"
)

func main() {
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	lg := logger.NewLogger("main : ", log.LstdFlags|log.Lshortfile)
	opt, _ := option.GetOptions()
	db, err := db.NewDB(ctx, opt.DatabaseURI)
	if err != nil {
		lg.Println(err)
		return
	}
	defer func() {
		err := db.Close()
		lg.Println(err)
	}()

	s, err := server.NewServer(db, opt.RunAddress)
	if err != nil{
		lg.Fatalln(err)
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
}
