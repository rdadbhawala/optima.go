package goroutine

import (
	"github.com/rdadbhawala/optima.go/optima"
)

// NewWorkshop is a goroutine based Workshop
func NewWorkshop(cfg *Config) optima.Workshop {
	w := &workshop{
		ch:      make(chan optima.Job),
		workers: make([]*worker, 0),
		c:       cfg,
	}
	w.AddWorker(cfg.Init)
	return w
}

type workshop struct {
	ch      chan optima.Job
	workers []*worker
	c       *Config
}

func (w *workshop) WorkerCount() int {
	return len(w.workers)
}

func (w *workshop) AddWorker(count int) error {
	if w.c.Max != 0 {
		count = w.minInt(count, w.c.Max-w.WorkerCount())
	}
	for i := 0; i < count; i++ {
		newW := newWorker(w.ch)
		w.workers = append(w.workers, newW)
		go newW.Start()
	}
	return nil
}

func (w *workshop) minInt(one, two int) int {
	if one < two {
		return one
	}
	return two
}

func (w *workshop) maxInt(one, two int) int {
	if one > two {
		return one
	}
	return two
}

func (w *workshop) RemoveWorker(count int) error {
	currLen := w.WorkerCount()
	if w.c.Min != 0 {
		count = w.minInt(count, currLen-w.c.Min)
	}
	newLen := currLen - count
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
