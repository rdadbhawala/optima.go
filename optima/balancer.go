package optima

// Balancer ...
type Balancer interface {
	Start()
}

// NewBalancer returns a new instance of a Balancer
func NewBalancer(pw Workshop, pjp JobProducer) Balancer {
	return &balancer{
		w:  pw,
		jp: pjp,
	}
}

type balancer struct {
	w  Workshop
	jp JobProducer
}

func (b *balancer) Start() {
	for {
		b.w.DoWork(b.jp.GetNextJob())
	}
}
