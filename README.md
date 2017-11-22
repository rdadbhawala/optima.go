# optima.go
Optima.go is a Go library for an in-memory, dynamic, self-optimizing worker pool.

## Story
Determining the right size of worker pool in a background service is not a difficult question to answer. It is however quite time-consuming as it is affected by several factors, many of which are external to the system.
* differing hardware capacity across various environments
* metrics of external dependencies such as database and web-service (and even the host OS) can be fluid
* (hardware + external dependencies + live load) is nearly impossible to simulate in a pre-live environment

To a certain extent, it may be possible to determine the ideal pool size after a series of experiments with the configurations. However, the influencing factors are not a constant, and hence the optimum determined based on certain criteria may not continue to be the right setting for a system.

The objective of Optima.go is to tackle the fluidity faced by a system (typically a background service) by adjusting its worker pool size.

## Concepts

**Workshop** is the pool of workers. You can add or remove workers from a Workshop. When there's a job to perform, it is passed to the Workshop, from where a worker thread can pick it up for processing.

**Producer** is the job factory. It produces jobs.

**Strategy** is where the worker pool is manipulated based on job completion events.

**Balancer** is the orchestration where Producer, Workshop and Strategy come together to deliver the promise of Optima.go.

## Implemenations

**optima/goroutine** is a Workshop based on Go-routines. Each worker runs as a go-routine. The workshop uses an unbuffered channel to pass jobs to available workers.

**optima/jobsPerSec** contains multiple implementations of the strategies based on jobs-per-second. The **Moving Average** strategy appears to be quite stable and reliable. It tracks a moving average of "jobs-per-second" for each batch size, and adjusts the pool according to the movement of this metric. There's also a **Simple Lever** strategy, prequel to Moving-Average, which checks the metrics for a batch of jobs based on the worker pool size.

**optima/basicBalancer** is a basic orchestration system.

Obviously, there's no implementation of **Producer**. That's your thing!

## Sample

### Getting optima.go
```
> go get github.com/rdadbhawala/optima.go/optima
```

### Using optima.go
```go
	// assuming p is the producer
	w := goroutine.NewWorkshop(&goroutine.Config{
		Min:  10,
		Max:  0,
		Init: 50,
	})
	s := jobsPerSec.NewMovingAverageStrategy(&jobsPerSec.MovingAverageConfig{
		Size:          5,
		WorkerRate:    25,
		PoolIncrement: 3,
	}, w)
	b := basicBalancer.NewBalancer(w, jp, s)
	go b.Start()
```
You can also try the testApp to look at the behavior of the algorithm on the console:
```
> go run .\testApp\main.go .\testApp\sleeper.go
```
It produces output which lists: index, worker pool size, batch average, new moving average, old moving average, diff, batch size.
The second last value indicates the key piece of logic: a positive diff increases the capacity, and a negative reduces it.
```
W: 0 130 1574.654 1787.9633 1875.0381 -87.07483 3325
W: 1 127 1797.7538 1751.6127 1787.9633 -36.350586 3250
W: 2 130 1779.918 1754.9178 1751.6127 3.3051758 3175
W: 3 133 1791.8135 1755.3043 1754.9178 0.3864746 3250
W: 4 130 1767.8662 1742.4006 1755.3043 -12.903687 3325
W: 0 133 2223.6362 1872.1971 1742.4006 129.79651 3250
W: 1 136 2236.2266 1959.8916 1872.1971 87.69446 3325
W: 2 139 2795.5337 2163.015 1959.8916 203.12329 3400
```
