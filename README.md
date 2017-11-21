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
