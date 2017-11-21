package jobsPerSec

import (
	"fmt"
	"time"

	"github.com/rdadbhawala/optima.go/optima"
)

// NewBalancer returns a new instance of a Balancer
func NewBalancer(pw optima.Workshop, pjp optima.JobProducer, pCfg Config) optima.Balancer {
	return &balancer{
		w:   pw,
		jp:  pjp,
		cfg: pCfg,
		ch:  make(chan *job),
	}
}

// Config ...
type Config struct {
	MeterHi       int
	MeterLo       int
	MeterInit     int
	ShakeThingsUp int
}

type balancer struct {
	cfg Config
	w   optima.Workshop
	jp  optima.JobProducer
	ch  chan *job
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
	prevTime := time.Now()
	meter := b.cfg.MeterInit
	cnt := 0
	prevJps := float32(0.0)
	sec := int(time.Second)
	unmodified := 0
	for {
		<-b.ch
		cnt++
		// b.val += j.end.Sub(j.start)
		if cnt >= b.w.WorkerCount()*25 {
			currTime := time.Now()
			dur := currTime.Sub(prevTime)
			newJps := float32(cnt*sec) / float32(dur)
			if newJps > prevJps {
				meter++
				if meter >= b.cfg.MeterHi {
					meter = 0
					unmodified = 0
					b.w.AddWorker(1)
				} else {
					unmodified++
				}
			} else {
				meter--
				if meter <= b.cfg.MeterLo {
					meter = 0
					unmodified = 0
					b.w.RemoveWorker(1)
				} else {
					unmodified++
				}
			}
			if unmodified >= b.cfg.ShakeThingsUp {
				meter = 0
				unmodified = 0
				b.w.AddWorker(1)
			}
			fmt.Println("W:", meter, b.w.WorkerCount(), prevJps, newJps, cnt)
			prevJps = newJps
			prevTime = currTime
			cnt = 0
		}
	}
}
