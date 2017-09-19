package jobsPerSec

import (
	"fmt"

	"github.com/rdadbhawala/optima.go/optima"
)

// NewBalancer returns a new instance of a Balancer
func NewBalancer(pw optima.Workshop, pjp optima.JobProducer) optima.Balancer {
	return &balancer{
		w:  pw,
		jp: pjp,
		ch: make(chan *job),
		sz: 100,
	}
}

type balancer struct {
	w       optima.Workshop
	jp      optima.JobProducer
	ch      chan *job
	sz      int
	cnt     int
	val     int64
	prevJps float32
}

func (b *balancer) Start() {
	go b.metrics()
	for {
		b.DoWork()
	}
}

func (b *balancer) DoWork() {
	j := newJob(b.jp.GetNextJob(), b.ch)
	b.w.DoWork(j)
}

func (b *balancer) metrics() {
	for {
		j := <-b.ch
		b.cnt++
		b.val += int64(j.end.Sub(j.start))
		if b.cnt == b.sz {
			newJps := float32(b.cnt) / float32(b.val)
			if newJps > b.prevJps {
				b.w.AddWorker(5)
				fmt.Println("Add:", b.w.WorkerCount(), b.prevJps, newJps)
				b.prevJps = newJps
			}
			b.val = 0
			b.cnt = 0
		}
	}
}
