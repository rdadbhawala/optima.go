package optima

// Job is a single unit of work.
type Job interface {
	DoWork()
}

// Producer is a source of work.
type Producer interface {
	GetNextJob() Job
}
