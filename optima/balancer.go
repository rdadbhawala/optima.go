package optima

// Balancer ...
type Balancer interface {
	Start()
}

// Strategy ...
type Strategy interface {
	JobCompleted(j Job)
}
