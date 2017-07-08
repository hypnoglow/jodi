package jodi

import "sync/atomic"

const (
	defaultMaxQueueSize = 2
)

// Dispatcher represents an offload service for job queue.
type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher.
	workerPool chan JobChan

	queue           JobChan
	jobErrorHandler JobErrorHandler
	maxWorkers      int
	numWaiting      int64
}

// NewDispatcher returns a new Dispatcher.
func NewDispatcher(maxWorkers int, jobErrorHandler JobErrorHandler) *Dispatcher {
	pool := make(chan JobChan, maxWorkers)
	queue := make(JobChan, defaultMaxQueueSize)
	return &Dispatcher{
		workerPool:      pool,
		queue:           queue,
		jobErrorHandler: jobErrorHandler,
		maxWorkers:      maxWorkers,
		numWaiting:      0,
	}
}

// Run creates a set of workers and starts listening
// for job requests from the in channel
// to dispatch them to the available worker.
func (d *Dispatcher) Run() {
	for i := 0; i < d.maxWorkers; i++ {
		// Each worker has the same pool so a worker can
		// write to it when it's ready to serve.
		NewWorker(d.workerPool, d.jobErrorHandler)
	}

	go d.dispatch()
}

// Enqueue adds a job to a queue. This will not block even if
// there is no available worker currently.
func (d *Dispatcher) Enqueue(job Job) {
	d.queue <- job
}

// NumWaiting returns a number of jobs waiting for an
// available worker.
func (d *Dispatcher) NumWaiting() int64 {
	return atomic.LoadInt64(&d.numWaiting)
}

func (d *Dispatcher) dispatch() {
	for job := range d.queue {
		// A job request has been received.

		// NOTE: This is a potential weakspot.
		// The amount of goroutines created here is not under control.
		// Another approach can be creating a local inmem queue instead
		// of running a separate goroutine for each pending job.
		go func(job Job) {
			atomic.AddInt64(&d.numWaiting, 1)

			availableWorkerChan := <-d.workerPool
			availableWorkerChan <- job

			atomic.AddInt64(&d.numWaiting, -1)
		}(job)
	}
}
