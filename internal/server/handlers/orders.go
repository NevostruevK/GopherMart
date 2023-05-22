package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/NevostruevK/GopherMart.git/internal/client"
	"github.com/NevostruevK/GopherMart.git/internal/db"
	"github.com/NevostruevK/GopherMart.git/internal/server/middleware"
	"github.com/NevostruevK/GopherMart.git/internal/util/luhn"
)

const (
	errWrongNumber                 = "неверный формат номера заказа"
	errExtractUserID               = `can't extract user ID`
	errOrderWasRegisterAnotherUser = "номер заказа уже был загружен другим пользователем"
	orderIsAccepted                = "новый номер заказа принят в обработку"
	errNoOrders                    = "нет данных для ответа"
)

func PostOrder(s *db.DB, m *client.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lg := middleware.GetLogger(r)
		lg.Println("PostOrder")
		b, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			lg.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		lg.Println(string(b))
		if !luhn.Valid(b) {
			lg.Println(errWrongNumber)
			http.Error(w, errWrongNumber, http.StatusUnprocessableEntity)
			return
		}
		userID, ok := r.Context().Value(middleware.KeyUserID).(uint64)
		if !ok {
			lg.Println(errExtractUserID)
			http.Error(w, errExtractUserID, http.StatusInternalServerError)
			return
		}
		orderUserID, err := s.PostOrder(r.Context(), userID, string(b))
		if err != nil {
			if !strings.Contains(err.Error(), db.ErrDuplicateOrder) {
				lg.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if orderUserID != userID {
				lg.Println(errOrderWasRegisterAnotherUser)
				http.Error(w, errOrderWasRegisterAnotherUser, http.StatusConflict)
				return
			}
			return
		}
		go m.NewTask(userID, string(b)).StandInLine()
		lg.Println(orderIsAccepted)
		w.WriteHeader(http.StatusAccepted)
		io.WriteString(w, orderIsAccepted)
	}
}

func GetOrders(s *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lg := middleware.GetLogger(r)
		lg.Println("GetOrders")
		_, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			lg.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		userID, ok := r.Context().Value(middleware.KeyUserID).(uint64)
		if !ok {
			lg.Println(errExtractUserID)
			http.Error(w, errExtractUserID, http.StatusInternalServerError)
			return
		}
		orders, err := s.GetOrders(r.Context(), userID)
		if err != nil {
			lg.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(orders) == 0 {
			lg.Println(errNoOrders)
			http.Error(w, errNoOrders, http.StatusNoContent)
			return
		}
		data, err := json.Marshal(&orders)
		if err != nil {
			lg.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}
