package basicBalancer

import "github.com/rdadbhawala/optima.go/optima"

// NewBalancer returns a new instance of a Balancer
func NewBalancer(ws optima.Workshop, jp optima.JobProducer, str optima.Strategy) optima.Balancer {
	return &basic{
		w:  ws,
		p:  jp,
		s:  str,
		ch: make(chan *job),
	}
}

type basic struct {
	w  optima.Workshop
	p  optima.JobProducer
	s  optima.Strategy
	ch chan *job
}

func (b *basic) Start() {
	go b.strategy()
	for {
		b.w.DoWork(newJob(b.p.GetNextJob(), b.ch))
	}
}

func (b *basic) strategy() {
	for {
		j := <-b.ch
		b.s.JobCompleted(j.j)
	}
}
