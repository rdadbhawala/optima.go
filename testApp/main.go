package main

import (
	"fmt"
	"time"

	"github.com/rdadbhawala/optima.go/optima/basicBalancer"

	"github.com/rdadbhawala/optima.go/optima/goroutine"
	"github.com/rdadbhawala/optima.go/optima/jobsPerSec"
)

func main() {
	jp := newSleepJobProducer(time.Millisecond * 150)
	w := goroutine.NewWorkshop(90)
	s := jobsPerSec.NewSimpleLeverStrategy(&jobsPerSec.SimpleLeverConfig{
		LeverHi:       2,
		LeverLo:       -2,
		LeverInit:     0,
		ShakeThingsUp: 10,
		WorkerRate:    25,
	}, w)
	b := basicBalancer.NewBalancer(w, jp, s)
	go b.Start()

	time.Sleep(1000 * time.Second)
	fmt.Println("JP", jp.count)
}
