package middleware

import (
	"context"
	"io"
	"net/http"

	"github.com/NevostruevK/GopherMart.git/internal/util/token"
)

const errNotAuthorization = "пользователь не аутентифицирован"

func AuthMiddleware(next http.Handler, tk *token.Token) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notAuth := []string{register, login}
		for _, path := range notAuth {
			if path == r.URL.Path {
				next.ServeHTTP(w, r)
				return
			}
		}
		lg := GetLogger(r)
		token := r.Header.Get("Authorization")
		ID, err := tk.Parse(token)
		if err != nil {
			lg.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			io.WriteString(w, errNotAuthorization)
			return
		}
		lg.Printf("Authorization user %d",ID)
		ctx := context.WithValue(r.Context(), KeyUserID, ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
