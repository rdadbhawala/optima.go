package basicBalancer

import (
	"github.com/rdadbhawala/optima.go/optima"
)

type job struct {
	j  optima.Job
	ch chan *job
}

func newJob(oj optima.Job, pch chan *job) *job {
	return &job{
		j:  oj,
		ch: pch,
	}
}

func (j *job) DoWork() {
	defer j.calc()
	j.j.DoWork()
}

func (j *job) calc() {
	j.ch <- j
}
