package jobsPerSec

import (
	"fmt"
	"time"

	"github.com/rdadbhawala/optima.go/optima"
)

// NewDefaultSimpleLeverConfig creates a default simple-lever config
func NewDefaultSimpleLeverConfig() *SimpleLeverConfig {
	return &SimpleLeverConfig{
		LeverHi:       2,
		LeverLo:       -2,
		LeverInit:     0,
		ShakeThingsUp: 10,
		WorkerRate:    25,
		PoolIncrement: 3,
	}
}

// SimpleLeverConfig is config for SimpleLeverStrategy
type SimpleLeverConfig struct {
	LeverHi       int
	LeverLo       int
	LeverInit     int
	ShakeThingsUp int
	WorkerRate    int
	PoolIncrement int
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
	if s.cnt >= s.w.WorkerCount()*s.c.WorkerRate {
		currTime := time.Now()
		dur := currTime.Sub(s.prevTime)
		newJps := float32(s.cnt*s.sec) / float32(dur)
		if newJps > s.prevJps {
			s.lever++
			if s.lever >= s.c.LeverHi {
				s.lever = s.c.LeverInit
				s.unmodified = 0
				s.w.AddWorker(s.c.PoolIncrement)
			} else {
				s.unmodified++
			}
		} else {
			s.lever--
			if s.lever <= s.c.LeverLo {
				s.lever = s.c.LeverInit
				s.unmodified = 0
				s.w.RemoveWorker(s.c.PoolIncrement)
			} else {
				s.unmodified++
			}
		}
		if s.unmodified >= s.c.ShakeThingsUp {
			s.lever = s.c.LeverInit
			s.unmodified = 0
			s.w.AddWorker(s.c.PoolIncrement)
		}
		fmt.Println("W:", s.lever, s.w.WorkerCount(), s.prevJps, newJps, s.cnt)
		s.prevJps = newJps
		s.prevTime = currTime
		s.cnt = 0
	}
}
