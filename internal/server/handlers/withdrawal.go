package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/NevostruevK/GopherMart.git/internal/db"
	"github.com/NevostruevK/GopherMart.git/internal/server/middleware"
	"github.com/NevostruevK/GopherMart.git/internal/util/luhn"
)

const (
	errSumUncorrect    = "не корректно задана сумма списания"
	errNotEnoughFounds = "на счету недостаточно средств"
	errNoWithdrawals   = "нет ни одного списания"
)

func PostWithdraw(s *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lg := middleware.GetLogger(r)
		lg.Println("PostWithdraw")
		b, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			lg.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		wOrder := db.WithdrawnOrder{}
		err = json.Unmarshal(b, &wOrder)
		if err != nil {
			lg.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		lg.Println(wOrder)
		if wOrder.Withdrawn == 0 {
			lg.Println(errSumUncorrect)
			http.Error(w, errSumUncorrect, http.StatusBadRequest)
			return
		}
		if !luhn.Valid([]byte(wOrder.Number)) {
			lg.Println(errWrongNumber)
			http.Error(w, errWrongNumber, http.StatusUnprocessableEntity)
			return
		}
		userID, ok := getUserID(w, r, lg)
		if !ok {
			return
		}
		if err = s.PostWithdrawal(r.Context(), userID, &wOrder); err != nil {
			if !strings.Contains(err.Error(), db.ErrNotEnoughFounds) {
				lg.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			lg.Println(errNotEnoughFounds)
			http.Error(w, errNotEnoughFounds, http.StatusPaymentRequired)
			return
		}
	}
}

func GetWithdrawals(s *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lg := middleware.GetLogger(r)
		lg.Println("GetWithdrawals")
		_, err := io.ReadAll(r.Body)
		if err != nil {
			lg.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		userID, ok := getUserID(w, r, lg)
		if !ok {
			return
		}
		wOrders, err := s.GetWithdrawals(r.Context(), userID)
		if err != nil {
			lg.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(wOrders) == 0 {
			lg.Println(errNoWithdrawals)
			http.Error(w, errNoWithdrawals, http.StatusNoContent)
			return
		}
		lg.Println(wOrders)
		data, err := json.Marshal(&wOrders)
		if err != nil {
			lg.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}
