package optima

// Strategy ...
type Strategy interface {
	JobCompleted(j Job)
}
