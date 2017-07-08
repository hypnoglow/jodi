package jodi

// Enqueuer can enqueue jobs.
type Enqueuer interface {
	// Enqueue adds a job to a queue.
	Enqueue(job Job)
}
