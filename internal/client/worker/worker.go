package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/NevostruevK/GopherMart.git/internal/client/task"
	"github.com/NevostruevK/GopherMart.git/internal/db"
	"github.com/NevostruevK/GopherMart.git/internal/util/fgzip"
	"github.com/NevostruevK/GopherMart.git/internal/util/logger"
)

type worker struct {
	client      *http.Client
	lg          *log.Logger
	url         string
	freeWorkers *int32
}

func NewWorker(address string, id uint64, free *int32) *worker {
	name := fmt.Sprintf("worker %d ", id)
	lg := logger.NewLogger(name, log.Lshortfile|log.LstdFlags)
	return &worker{&http.Client{}, lg, address, free}
}

func (w worker) free() {
	atomic.AddInt32(w.freeWorkers, 1)
}

func (w worker) busy() {
	atomic.AddInt32(w.freeWorkers, -1)
}

func (w worker) Start(ctx context.Context, s *db.DB, ch chan task.Task) {
	w.lg.Println("Start")
	w.free()
	for {
		select {
		case task := <-ch:
			w.lg.Printf("Get task %d", task.TaskID)
			w.busy()
			order, code, err := w.getOrder(ctx, task.Order.Number)
			if err != nil || code != http.StatusOK {
				w.wrongCompletition(ctx, task, err, code)
				break
			}
			if ok := task.NeedUpdateOrder(*order); ok {
				err = s.UpdateOrder(ctx, task.UserID, task.Order)
				if err != nil {
					w.wrongCompletition(ctx, task, err, code)
					break
				}
				w.lg.Printf("task %d complete", task.TaskID)
				task.Finished = true
				go task.StandInLine()
				w.free()
			}
		case <-ctx.Done():
			w.lg.Println("Finished")
			w.free()
			return
		}
	}
}
func (w worker) wrongCompletition(ctx context.Context, task task.Task, err error, code int) {
	if err != nil {
		w.lg.Printf("task %d complete with err %v", task.TaskID, err)
		task.SetError(err)
	} else { // code != http.StatusOK
		w.lg.Printf("task %d complete with code %d", task.TaskID, code)
		task.SetTry(code)
	}
	go task.StandInLine()
	w.free()
}

func (w worker) getOrder(ctx context.Context, number string) (*task.Order, int, error) {
	w.url += "/api/orders/"+number
	w.lg.Println(w.url)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, w.url, nil)
	if err != nil {
		w.lg.Println(err)
		return nil, 0, err
	}
	response, err := w.client.Do(request)
	if err != nil {
		w.lg.Println(err)
		return nil, 0, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		w.lg.Println(err)
		return nil, 0, err
	}
	if response.StatusCode != http.StatusOK {
		w.lg.Printf("Status code %d", response.StatusCode)
		return nil, response.StatusCode, nil
	}
	if strings.Contains(response.Header.Get("Content-Encoding"), "gzip") {
		body, err = fgzip.Decompress(body)
		if err != nil {
			w.lg.Println(err)
			return nil, http.StatusOK, err
		}
	}
	order := task.Order{}
	err = json.Unmarshal(body, &order)
	if err != nil {
		w.lg.Println(err)
		return nil, http.StatusOK, err
	}
	return &order, http.StatusOK, nil
}
