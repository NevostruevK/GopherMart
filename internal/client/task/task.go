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

func (o Order) String() string {
	if o.Accrual == nil {
		return fmt.Sprintf("%s : %s ", o.Number, o.Status)
	}
	return fmt.Sprintf("%s : %s : %f", o.Number, o.Status, *o.Accrual)
}

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
		fifo:         fifo,
		Order:        &Order{Number: number, Status: NEW},
	}
}
func (t Task) StandInLine() {
	t.fifo <- t
}

func (t *Task) SetError(err error) {
	t.ErrorsNumber++
}

func (t *Task) SetTry(code int) {
	t.TriesNumber++
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
