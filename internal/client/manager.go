package client

import (
	"context"
	"log"
	"sync/atomic"

	"github.com/NevostruevK/GopherMart.git/internal/client/task"
	"github.com/NevostruevK/GopherMart.git/internal/db"
	"github.com/NevostruevK/GopherMart.git/internal/util/logger"
)

const taskInitialCount = 16

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

func (m *Manager) Start(ctx context.Context, s *db.DB, address string, workersCount int32) {
	lg := logger.NewLogger("manager : ", log.LstdFlags|log.Lshortfile)
	lg.Println("Start")
	//	chOut := make(chan task.Task)
	tasks := make([]task.Task, 0, taskInitialCount)
	for {
		select {
		case newTask := <-m.inCh:
			tasks = append(tasks, newTask)
			lg.Println("take a task")
		case <-ctx.Done():
			lg.Println("Finished")
			return
		}
	}
}
