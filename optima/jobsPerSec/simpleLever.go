package jobsPerSec

import (
	"fmt"
	"time"

	"github.com/rdadbhawala/optima.go/optima"
)

// SimpleLeverConfig is config for SimpleLeverStrategy
type SimpleLeverConfig struct {
	LeverHi       int
	LeverLo       int
	LeverInit     int
	ShakeThingsUp int
	WorkerRate    int
}

// NewSimpleLeverStrategy returns SimpleLeverStrategy
func NewSimpleLeverStrategy(cf *SimpleLeverConfig, ws optima.Workshop) optima.Strategy {
	return &simpleLeverStrategy{
		c:          cf,
		w:          ws,
		prevTime:   time.Now(),
		lever:      cf.LeverInit,
		cnt:        0,
		prevJps:    float32(0.0),
		sec:        int(time.Second),
		unmodified: 0,
	}
}

type simpleLeverStrategy struct {
	c          *SimpleLeverConfig
	w          optima.Workshop
	prevTime   time.Time
	lever      int
	cnt        int
	prevJps    float32
	sec        int
	unmodified int
}

func (s *simpleLeverStrategy) JobCompleted(j optima.Job) {
	s.cnt++
	// b.val += j.end.Sub(j.start)
	if s.cnt >= s.w.WorkerCount()*s.c.WorkerRate {
		currTime := time.Now()
		dur := currTime.Sub(s.prevTime)
		newJps := float32(s.cnt*s.sec) / float32(dur)
		if newJps > s.prevJps {
			s.lever++
			if s.lever >= s.c.LeverHi {
				s.lever = 0
				s.unmodified = 0
				s.w.AddWorker(1)
			} else {
				s.unmodified++
			}
		} else {
			s.lever--
			if s.lever <= s.c.LeverLo {
				s.lever = 0
				s.unmodified = 0
				s.w.RemoveWorker(1)
			} else {
				s.unmodified++
			}
		}
		if s.unmodified >= s.c.ShakeThingsUp {
			s.lever = 0
			s.unmodified = 0
			s.w.AddWorker(1)
		}
		fmt.Println("W:", s.lever, s.w.WorkerCount(), s.prevJps, newJps, s.cnt)
		s.prevJps = newJps
		s.prevTime = currTime
		s.cnt = 0
	}
}
