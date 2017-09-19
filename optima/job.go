package optima

// Job is a single unit of work.
type Job interface {
	DoWork()
}

// JobProducer is a source of work.
type JobProducer interface {
	GetNextJob() Job
}
