package jobsPerSec

import (
	"fmt"
	"time"

	"github.com/rdadbhawala/optima.go/optima"
)

// NewBalancer returns a new instance of a Balancer
func NewBalancer(pw optima.Workshop, pjp optima.JobProducer) optima.Balancer {
	return &balancer{
		w:        pw,
		jp:       pjp,
		ch:       make(chan *job),
		sz:       10,
		meterHi:  2,
		meterLo:  -2,
		prevTime: time.Now().UTC(),
	}
}

type balancer struct {
	w        optima.Workshop
	jp       optima.JobProducer
	ch       chan *job
	sz       int
	cnt      int
	val      time.Duration
	prevJps  float32
	meter    int
	meterHi  int
	meterLo  int
	prevTime time.Time
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
	sec := int(time.Second)
	unmodified := 0
	unmodifiedRange := (b.meterHi - b.meterLo + 1) * 2
	for {
		<-b.ch
		b.cnt++
		// b.val += j.end.Sub(j.start)
		if b.cnt >= b.w.WorkerCount()*25 {
			currTime := time.Now().UTC()
			dur := currTime.Sub(b.prevTime)
			newJps := float32(b.cnt*sec) / float32(dur)
			if newJps > b.prevJps {
				b.meter++
				if b.meter >= b.meterHi {
					b.meter = 0
					unmodified = 0
					b.w.AddWorker(1)
				} else {
					unmodified++
				}
			} else {
				b.meter--
				if b.meter <= b.meterLo {
					b.meter = 0
					unmodified = 0
					b.w.RemoveWorker(1)
				} else {
					unmodified++
				}
			}
			if unmodified >= unmodifiedRange {
				b.meter = 0
				unmodified = 0
				b.w.AddWorker(1)
			}
			fmt.Println("W:", b.meter, b.w.WorkerCount(), b.prevJps, newJps, b.cnt)
			b.prevJps = newJps
			b.prevTime = currTime
			b.val = 0
			b.cnt = 0
		}
	}
}
