package jobsPerSec

import (
	"time"

	"github.com/rdadbhawala/optima.go/optima"
)

type job struct {
	j     optima.Job
	start time.Time
	end   time.Time
	ch    chan *job
}

func newJob(oj optima.Job, pch chan *job) *job {
	return &job{
		j:  oj,
		ch: pch,
	}
}

func (j *job) DoWork() {
	defer j.calc()
	j.start = time.Now()
	j.j.DoWork()
}

func (j *job) calc() {
	j.end = time.Now()
	j.ch <- j
}
