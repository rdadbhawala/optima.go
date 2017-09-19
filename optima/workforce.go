package optima

// Workshop is a set of Workers.
type Workshop interface {
	WorkerCount() int
	AddWorker(count int) error
	RemoveWorker(count int) error
	DoWork(j Job)
}

// Worker does some work.
type Worker interface {
}
