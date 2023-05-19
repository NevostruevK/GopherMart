package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"

	"github.com/NevostruevK/GopherMart.git/internal/client/task"
	"github.com/NevostruevK/GopherMart.git/internal/util/fgzip"
	"github.com/NevostruevK/GopherMart.git/internal/util/logger"
)

type worker struct {
	client *http.Client
	lg     *log.Logger
	url    url.URL
}

func NewWorker(address string, id uint64) *worker {
	name := fmt.Sprintf("worker %d ", id)
	lg := logger.NewLogger(name, log.Lshortfile|log.LstdFlags)
	url := url.URL{
		Scheme: "http",
		Host:   address + "/",
	}
	return &worker{&http.Client{}, lg, url}
}
func (w worker) Start(ctx context.Context, ch chan task.Task, free *int32) {
	w.lg.Println("Start")
	atomic.AddInt32(free, 1)
	for {
		select {
		case task := <-ch:
			w.lg.Printf("Get task %d", task.TaskID)
			atomic.AddInt32(free, -1)
			order, code, err := w.getOrder(ctx, task.Order.Number)
			if err != nil {
				w.lg.Println(err)
				task.SetError(err)
				go task.StandInLine()
				atomic.AddInt32(free, 1)
				break
			}
			if code != http.StatusOK {
				w.lg.Println(err)
				task.SetTry(code)
				go task.StandInLine()
				atomic.AddInt32(free, 1)
				break
			}
			if ok := task.NeedUpdateOrder(*order); ok {
				goDatatBase
			}
		case <-ctx.Done():
			w.lg.Println("Finished")
			atomic.AddInt32(free, 1)
			return
		}
	}
}

func (w worker) getOrder(ctx context.Context, number string) (*task.Order, int, error) {
	w.url.Path = number
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, w.url.String(), nil)
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
