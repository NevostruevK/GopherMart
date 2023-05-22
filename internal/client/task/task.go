package task

import (
	"fmt"
)

type Order struct {
	Number  string   `json:"order"`             // Номер заказа
	Status  status   `json:"status"`            // Статус заказа
	Accrual *float64 `json:"accrual,omitempty"` // Начислено баллов
}
type status string

const (
	NEW        status = "NEW"
	REGISTERED status = "REGISTERED"
	INVALID    status = "INVALID"
	PROCESSING status = "PROCESSING"
	PROCESSED  status = "PROCESSED"
)

const (
	errOrderAccrualNil      = "accrual is nil"
	errOrderAccrualNegative = "accrual is negative"
)

func (o Order) Valid() error {
	if o.Status != PROCESSED {
		return nil
	}
	if o.Accrual == nil {
		return fmt.Errorf(errOrderAccrualNil)
	}
	if *o.Accrual < 0 {
		return fmt.Errorf(errOrderAccrualNegative)
	}
	return nil
}

type Task struct {
	UserID       uint64
	TaskID       uint64
	ErrorsNumber uint32
	TriesNumber  uint32
	Finished     bool
//	lg           log.Logger
	fifo         chan Task
	Order        *Order
}

func NewTask(userID, taskID uint64, fifo chan Task, number string) *Task {
	return &Task{
		UserID:       userID,
		TaskID:       taskID,
		ErrorsNumber: 0,
		TriesNumber:  0,
		Finished:     false,
//		lg: *logger.NewLogger(fmt.Sprintf("task %d :", taskID),
//			log.Lshortfile|log.LstdFlags),
		fifo:  fifo,
		Order: &Order{Number: number, Status: NEW},
	}
}
func (t Task) StandInLine() {
//	t.lg.Println("stand in line")
	t.fifo <- t
//	t.lg.Println("manager took me")
}

func (t *Task) SetError(err error) {
	t.ErrorsNumber++
//	t.lg.Printf("ERROR : %v", err)
}

func (t *Task) SetTry(code int) {
	t.TriesNumber++
//	t.lg.Printf("CODE : %d", code)
}

func (t *Task) NeedUpdateOrder(o Order) bool {
	if err := o.Valid(); err != nil {
		t.SetError(err)
		return false
	}
	if t.Order.Status == o.Status {
		return false
	}
	t.Order.Status = o.Status
	t.Order.Accrual = o.Accrual
	return true
}
