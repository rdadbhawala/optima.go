package jobsPerSec

import (
	"fmt"
	"time"

	"github.com/rdadbhawala/optima.go/optima"
)

// MovingAverageConfig is a config
type MovingAverageConfig struct {
	Size          int
	WorkerRate    int
	PoolIncrement int
}

// NewMovingAverageStrategy returns a movingAverageStrategy with configuration
func NewMovingAverageStrategy(cf *MovingAverageConfig, ws optima.Workshop) optima.Strategy {
	s := &movingAverageStrategy{
		c:          cf,
		w:          ws,
		cnt:        0,
		averages:   make([]float32, cf.Size),
		sum:        0.0,
		prevTime:   time.Now(),
		index:      0,
		sec:        int(time.Second),
		sizeFactor: 1.0 / float32(cf.Size),
		prevAvg:    0.0,
	}
	for ctr := 0; ctr < cf.Size; ctr++ {
		s.averages[ctr] = 0
	}
	return s
}

type movingAverageStrategy struct {
	c          *MovingAverageConfig
	w          optima.Workshop
	cnt        int
	averages   []float32
	sum        float32
	prevTime   time.Time
	index      int
	sec        int
	sizeFactor float32
	prevAvg    float32
}

func (s *movingAverageStrategy) JobCompleted(j optima.Job) {
	s.cnt++
	if s.cnt >= s.c.WorkerRate*s.w.WorkerCount() {
		currTime := time.Now()
		dur := currTime.Sub(s.prevTime)
		newJps := float32(s.cnt*s.sec) / float32(dur)
		s.sum += (newJps - s.averages[s.index])
		newAvg := s.sum * s.sizeFactor
		if newAvg > s.prevAvg {
			s.w.AddWorker(s.c.PoolIncrement)
		} else {
			s.w.RemoveWorker(s.c.PoolIncrement)
		}
		fmt.Println("W:", s.index, s.w.WorkerCount(), newJps, newAvg, s.prevAvg, newAvg-s.prevAvg, s.cnt)
		s.averages[s.index] = newJps
		s.prevAvg = newAvg
		s.cnt = 0
		s.index = (s.index + 1) % s.c.Size
		s.prevTime = currTime
	}
}
