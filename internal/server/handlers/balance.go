package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/NevostruevK/GopherMart.git/internal/db"
	"github.com/NevostruevK/GopherMart.git/internal/server/middleware"
)

func GetBalance(s *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lg := middleware.GetLogger(r)
		lg.Println("GetBalance")
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
		balance, err := s.GetBalance(r.Context(), userID)
		if err != nil {
			lg.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		lg.Println(balance)
		data, err := json.Marshal(&balance)
		if err != nil {
			lg.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}
