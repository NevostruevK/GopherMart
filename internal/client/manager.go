package client

import (
	"context"
	"log"
	"sync/atomic"

	"github.com/NevostruevK/GopherMart.git/internal/client/task"
	"github.com/NevostruevK/GopherMart.git/internal/client/worker"
	"github.com/NevostruevK/GopherMart.git/internal/db"
	"github.com/NevostruevK/GopherMart.git/internal/util/logger"
)

const taskInitialCount = 16

const (
	needErrorsToReject = 32
	needTriesToReject  = 32
)

type Manager struct {
	taskID      uint64
	freeWorkers int32
	inCh        chan task.Task
}

func NewManager() *Manager {
	return &Manager{taskID: 0, inCh: make(chan task.Task)}
}

func (m *Manager) NewTask(userID uint64, number string) *task.Task {
	return task.NewTask(userID, atomic.AddUint64(&m.taskID, 1), m.inCh, number)
}

func (m *Manager) Start(ctx context.Context, s *db.DB, address string, workersCount int) {
	var completed uint64
	var rejected uint64
	var took uint64
	lg := logger.NewLogger("manager : ", log.LstdFlags|log.Lshortfile)
	lg.Println("Start")
	chOut := make(chan task.Task)
	for i := 0; i < workersCount; i++ {
		w := worker.NewWorker(address, uint64(i), &m.freeWorkers)
		go w.Start(ctx, s, chOut)
	}
	tasks := make([]task.Task, 0, taskInitialCount)
	for {
		select {
		case task := <-m.inCh:
			if task.ErrorsNumber >= needErrorsToReject || task.TriesNumber >= needTriesToReject {
				lg.Printf("rejected task %d", task.TaskID)
				rejected++
				break
			}
			if task.Finished {
				lg.Printf("completed task %d", task.TaskID)
				completed++
				break
			}
			tasks = append(tasks, task)
			took++
			lg.Printf("take task %d", task.TaskID)
		case <-ctx.Done():
			lg.Println("Finished")
			lg.Printf("took tasks %d", took)
			lg.Printf("completed tasks %d", completed)
			lg.Printf("rejected tasks %d", rejected)
			return
		default:
			if len(tasks) > 0 && atomic.LoadInt32(&m.freeWorkers) > 0 {
				chOut <- tasks[0]
				lg.Printf("Gave task %d", tasks[0].TaskID)
				tasks = tasks[1:]
			}
		}
	}
}
