package main

import (
	"fmt"
	"time"

	"github.com/rdadbhawala/optima.go/optima"
	"github.com/rdadbhawala/optima.go/optima/goroutine"
)

func main() {
	jp := newSleepJobProducer(time.Millisecond * 300)
	b := optima.NewBalancer(goroutine.NewWorkshop(), jp)
	go b.Start()

	time.Sleep(5 * time.Second)
	fmt.Println("JP", jp.count)
}
