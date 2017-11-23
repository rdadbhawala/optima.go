package optima

// Job is a single unit of work.
type Job interface {
	DoWork()
}

// Producer is a source of work.
type Producer interface {
	GetNextJob() Job
}

// DoWork is a function, useful when job types are inline in code
type DoWork func()

// DoWorkInJob wraps a DoWork function in an optima.Job so that it can be used by the system
func DoWorkInJob(d DoWork) Job {
	return &doWorkJob{
		dw: d,
	}
}

type doWorkJob struct {
	dw DoWork
}

func (d *doWorkJob) DoWork() {
	d.dw()
}
