package main

import (
	"sync"
	"time"

	"github.com/rdadbhawala/optima.go/optima"
)

// newSleepJobProducer ... returns an instance of Sleeper Job Producer
func newSleepJobProducer(t time.Duration) *sleepJobProducer {
	return &sleepJobProducer{
		sleepTime: t,
	}
}

// sleepJobProducer is the Job Producer
type sleepJobProducer struct {
	count     int
	lock      sync.Mutex
	sleepTime time.Duration
}

// GetNextJob gets the next job
func (s *sleepJobProducer) GetNextJob() optima.Job {
	defer s.lock.Unlock()
	s.lock.Lock()
	s.count++
	return &sleepJob{
		index:     s.count,
		sleepTime: s.sleepTime,
	}
}

// sleepJob is a Job
type sleepJob struct {
	index     int
	sleepTime time.Duration
}

// DoWork does the work
func (j *sleepJob) DoWork() {
	// fmt.Println("Job: ", j.index, j.sleepTime)
	limit := (j.sleepTime + time.Duration(j.index%100)) / 4
	// for cnt := 0; cnt < 4; cnt++ {
	time.Sleep(limit)
	for i := 0; i < int(limit/100); i++ {
	}
	// }
}
