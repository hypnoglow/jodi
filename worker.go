package jodi

// Worker represents the worker that runs the Job.
type Worker struct {
	pool            chan<- JobChan
	jobErrorHandler JobErrorHandler
	jobChan         chan Job
	quit            chan empty
}

// NewWorker creates a new worker,
// starts listening for job requests
// and returns a Worker handle.
func NewWorker(pool chan<- JobChan, jobErrorHandler JobErrorHandler) Worker {
	w := Worker{
		pool:            pool,
		jobErrorHandler: jobErrorHandler,
		jobChan:         make(chan Job),
		quit:            make(chan empty),
	}

	go w.start()
	return w
}

// start signals the worker to start listening for job requests.
func (w Worker) start() {
	go func() {
		for {
			// Register the current worker in the pool.
			w.pool <- w.jobChan

			select {
			case job := <-w.jobChan:
				if err := job.Run(); err != nil {
					w.jobErrorHandler(err)
				}
			case <-w.quit:
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for job requests.
func (w Worker) Stop() {
	close(w.quit)
}

type empty struct{}
