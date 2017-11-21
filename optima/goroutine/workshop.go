package goroutine

import (
	"github.com/rdadbhawala/optima.go/optima"
)

// NewWorkshop is a goroutine based Workshop
func NewWorkshop(workerCount int) optima.Workshop {
	w := &workshop{
		ch:      make(chan optima.Job),
		workers: make([]*worker, 0),
		min:     workerCount,
	}
	w.AddWorker(w.min)
	return w
}

type workshop struct {
	jp      optima.Producer
	ch      chan optima.Job
	workers []*worker
	min     int
}

func (w *workshop) WorkerCount() int {
	return len(w.workers)
}

func (w *workshop) AddWorker(count int) error {
	for i := 0; i < count; i++ {
		newW := newWorker(w.ch)
		w.workers = append(w.workers, newW)
		go newW.Start()
	}
	return nil
}

func (w *workshop) RemoveWorker(count int) error {
	currLen := w.WorkerCount()
	newLen := currLen - count
	if w.min > newLen {
		newLen = w.min
	}
	dropWorkers := w.workers[newLen:]
	for cnt := 0; cnt < len(dropWorkers); cnt++ {
		dropWorkers[cnt].Stop()
	}
	w.workers = w.workers[:newLen]
	return nil
}

func (w *workshop) DoWork(j optima.Job) {
	w.ch <- j
}
