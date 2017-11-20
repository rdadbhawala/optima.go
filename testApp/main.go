package main

import (
	"fmt"
	"time"

	"github.com/rdadbhawala/optima.go/optima/goroutine"
	"github.com/rdadbhawala/optima.go/optima/jobsPerSec"
)

func main() {
	jp := newSleepJobProducer(time.Millisecond * 150)
	w := goroutine.NewWorkshop(1)
	//w.AddWorker(4)
	b := jobsPerSec.NewBalancer(w, jp)
	go b.Start()

	time.Sleep(1000 * time.Second)
	fmt.Println("JP", jp.count)
}
