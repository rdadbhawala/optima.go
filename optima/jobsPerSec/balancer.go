package jobsPerSec

import (
	"fmt"

	"github.com/rdadbhawala/optima.go/optima"
)

// NewBalancer returns a new instance of a Balancer
func NewBalancer(pw optima.Workshop, pjp optima.JobProducer) optima.Balancer {
	return &balancer{
		w:       pw,
		jp:      pjp,
		ch:      make(chan *job),
		sz:      10,
		meterHi: 2,
		meterLo: -2,
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
	meter   int
	meterHi int
	meterLo int
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
		if b.cnt == (b.sz * b.w.WorkerCount()) {
			newJps := float32(b.cnt) / float32(b.val)
			if newJps > b.prevJps {
				if b.meter < b.meterHi {
					b.meter++
				}
				if b.meter == b.meterHi {
					// b.meter = 0
					b.w.AddWorker(5)
				}
			} else {
				if b.meter > b.meterLo {
					b.meter--
				}
				if b.meter == b.meterLo {
					// b.meter = 0
					b.w.RemoveWorker(5)
				}
			}
			fmt.Println("W:", b.meter, b.w.WorkerCount(), b.prevJps, newJps)
			b.prevJps = newJps
			b.val = 0
			b.cnt = 0
		}
	}
}
