package middleware

import (
	"io"
	"net/http"
	"strings"
)

const (
	errNotJson = "Content-Type is not application/json"
	errNotText = "Content-Type is not text/plain"
)
const (
	register       = "/api/user/register"
	login          = "/api/user/login"
	postOrder      = "/api/user/orders"
	postWithdrawal = "/api/user/balance/withdraw"
)

func CheckHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			next.ServeHTTP(w, r)
			return
		}
		lg := GetLogger(r)
		switch r.URL.Path {
		case postOrder:
			if !strings.Contains(r.Header.Get("Content-Type"), "text/plain") {
				lg.Println(errNotText)
				w.WriteHeader(http.StatusBadRequest)
				io.WriteString(w, errNotText)
				return
			}
		default:
			if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
				lg.Println(errNotJson)
				w.WriteHeader(http.StatusBadRequest)
				io.WriteString(w, errNotJson)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
