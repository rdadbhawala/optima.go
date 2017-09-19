package goroutine

import "github.com/rdadbhawala/optima.go/optima"

func newWorker(pch chan optima.Job) *worker {
	return &worker{
		ch:  pch,
		run: true,
	}
}

type worker struct {
	ch  chan optima.Job
	run bool
}

func (w *worker) DoWork() {
	j := <-w.ch
	j.DoWork()
}

func (w *worker) Start() {
	for w.run {
		w.DoWork()
	}
}

func (w *worker) Stop() {
	w.run = false
}
