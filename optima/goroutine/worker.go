package goroutine

import "github.com/rdadbhawala/optima.go/optima"

func newWorker(pch chan optima.Job) *worker {
	return &worker{
		ch: pch,
	}
}

type worker struct {
	ch chan optima.Job
}

func (w *worker) DoWork() {
	j := <-w.ch
	j.DoWork()
}

func (w *worker) Start() {
	for {
		w.DoWork()
	}
}
