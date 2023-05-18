package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/NevostruevK/GopherMart.git/internal/db"
	"github.com/NevostruevK/GopherMart.git/internal/server/middleware"
	"github.com/NevostruevK/GopherMart.git/internal/util/token"
)

const (
	errLoginAlreadyOccupied   = "логин уже занят"
	errLoginPasswordUncorrect = "не корректно заданы логин пароль"
	errWrongLoginPassword     = "неверная пара логин/пароль"
)

type requestType int

const (
	Register requestType = iota
	Login
)

func (r requestType) String() string {
	if r == Register {
		return "Register"
	}
	return "Login"
}
func (r requestType) goDataBase(ctx context.Context, s *db.DB, u db.User) (uint64, error) {
	if r == Register {
		return s.Register(ctx, u)
	}
	return s.Login(ctx, u)
}

func Authentication(s *db.DB, tk *token.Token, rType requestType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lg := middleware.GetLogger(r)
		lg.Println(rType)
		b, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			lg.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		u := db.User{}
		err = json.Unmarshal(b, &u)
		if err != nil {
			lg.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		lg.Println(u)
		if u.Login == "" || u.Password == "" {
			lg.Println(errLoginPasswordUncorrect)
			http.Error(w, errLoginPasswordUncorrect, http.StatusBadRequest)
			return
		}
		ID, err := rType.goDataBase(r.Context(), s, u)
		if err != nil {
			if strings.Contains(err.Error(), db.ErrDuplicateLogin) {
				if rType == Register {
					lg.Println(errLoginAlreadyOccupied)
					http.Error(w, errLoginAlreadyOccupied, http.StatusConflict)
					return
				}
				lg.Println(errWrongLoginPassword)
				http.Error(w, errWrongLoginPassword, http.StatusUnauthorized)
				return
			}
			lg.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		token, err := tk.Make(ID)
		if err != nil {
			lg.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		lg.Printf("login user %d",ID)
		w.Header().Set("Authorization", token)

		//		io.WriteString(w, fmt.Sprintln(id))
		//	w.Write(id)
	}
}
