package jobsPerSec

import (
	"fmt"
	"time"

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
		// prevTime: time.Now().UTC(),
	}
}

type balancer struct {
	w       optima.Workshop
	jp      optima.JobProducer
	ch      chan *job
	sz      int
	cnt     int
	val     time.Duration
	prevJps float32
	meter   int
	meterHi int
	meterLo int
	// prevTime time.Time
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
	sec := float32(time.Second)
	unmodified := 0
	for {
		j := <-b.ch
		b.cnt++
		b.val += j.end.Sub(j.start)
		if b.cnt >= 100 {
			// currTime := time.Now().UTC()
			newSpj := float32(b.val) / float32(b.cnt) / sec
			newJps := float32(b.cnt) / float32(b.val) * sec
			fmt.Println("W:", b.meter, b.w.WorkerCount(), newJps, b.cnt, newSpj, b.val)
			if newJps > b.prevJps {
				if b.meter < b.meterHi {
					b.meter++
				}
				if b.meter == b.meterHi {
					b.meter = 0
					unmodified = 0
					b.w.AddWorker(1)
				} else {
					unmodified++
				}
			} else {
				if b.meter > b.meterLo {
					b.meter--
				}
				if b.meter == b.meterLo {
					b.meter = 0
					unmodified = 0
					b.w.RemoveWorker(1)
				} else {
					unmodified++
				}
			}
			if unmodified >= 5 {
				b.meter = 0
				unmodified = 0
				b.w.AddWorker(1)
			}
			b.prevJps = newJps
			b.val = 0
			b.cnt = 0
		}
	}
}
