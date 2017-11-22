package main

import (
	"fmt"
	"time"

	"github.com/rdadbhawala/optima.go/optima/basicBalancer"

	"github.com/rdadbhawala/optima.go/optima/goroutine"
	"github.com/rdadbhawala/optima.go/optima/jobsPerSec"
)

func main() {
	movingAverage()
}

func movingAverage() {
	jp := newSleepJobProducer(time.Millisecond * 150)
	w := goroutine.NewWorkshop(&goroutine.Config{
		Min:  10,
		Max:  0,
		Init: 100,
	})
	s := jobsPerSec.NewMovingAverageStrategy(&jobsPerSec.MovingAverageConfig{
		Size:          5,
		WorkerRate:    25,
		PoolIncrement: 3,
	}, w)
	b := basicBalancer.NewBalancer(w, jp, s)
	go b.Start()

	time.Sleep(1000 * time.Second)
	fmt.Println("JP", jp.count)
}

func simpleLever() {
	jp := newSleepJobProducer(time.Millisecond * 150)
	w := goroutine.NewWorkshop(&goroutine.Config{
		Min:  10,
		Max:  0,
		Init: 100,
	})
	s := jobsPerSec.NewSimpleLeverStrategy(&jobsPerSec.SimpleLeverConfig{
		LeverHi:       3,
		LeverLo:       -3,
		LeverInit:     0,
		ShakeThingsUp: 12,
		WorkerRate:    25,
		PoolIncrement: 3,
	}, w)
	b := basicBalancer.NewBalancer(w, jp, s)
	go b.Start()

	time.Sleep(1000 * time.Second)
	fmt.Println("JP", jp.count)
}
