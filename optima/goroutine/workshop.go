package goroutine

import "github.com/rdadbhawala/optima.go/optima"

// NewWorkshop is a goroutine based Workshop
func NewWorkshop() optima.Workshop {
	w := &workshop{
		ch:      make(chan optima.Job),
		workers: make([]*worker, 0),
	}
	w.AddWorker(1)
	return w
}

type workshop struct {
	jp      optima.JobProducer
	ch      chan optima.Job
	workers []*worker
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
	return nil
}

func (w *workshop) DoWork(j optima.Job) {
	w.ch <- j
}
