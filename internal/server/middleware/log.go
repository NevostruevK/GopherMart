package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/NevostruevK/GopherMart.git/internal/util/logger"
)

var requestID uint64 = 0

type ctxKey string

const (
	keyServerLog ctxKey = "log"
	KeyUserID    ctxKey = "userID"
)

func newServerLogger() *log.Logger {
	ID := atomic.AddUint64(&requestID, 1)
	return logger.NewLogger(fmt.Sprintf("request %d : ", ID), log.Lshortfile|log.LstdFlags)
}

func GetLogger(r *http.Request) *log.Logger {
	lg, ok := r.Context().Value(keyServerLog).(*log.Logger)
	if ok {
		return lg
	}
	return logger.NewLogger("default server logger", log.Lshortfile|log.LstdFlags)
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lg := newServerLogger()
		lg.Printf("URL : %s", r.URL)
		ctx := context.WithValue(r.Context(), keyServerLog, lg)
		next.ServeHTTP(w, r.WithContext(ctx))
		lg.Println("complete")
	})
}
