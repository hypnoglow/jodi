package jodi

// Job describes the job to be run.
type Job interface {
	// Run runs the job.
	Run() error
}

// JobChan describes a channel that we can send jobs to and receive jobs from.
type JobChan chan Job

// JobErrorHandler describes a func that is used to handle job error.
type JobErrorHandler func(error)
