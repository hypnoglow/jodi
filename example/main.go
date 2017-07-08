package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hypnoglow/jodi"
)

const (
	defaultMaxWorkers = 4
	numberOfJobs      = 50

	dispatcherMonitoringPeriod = time.Millisecond * 500
	senderPeriod               = time.Millisecond * 200
)

func main() {
	rand.Seed(time.Now().Unix())

	// Prepare a function to handle errors in jobs.
	errorHandler := func(err error) {
		log.Fatalln("Job failed: %s", err.Error())
	}

	// Create a new Dispatcher and Run it to start delegating
	// jobs to workers.
	dp := jodi.NewDispatcher(defaultMaxWorkers, errorHandler)
	dp.Run()

	// Monitoring number of jobs waiting for an available worker.
	go func() {
		for {
			log.Println("[DSPTCHR] Number of waiting: ", dp.NumWaiting())
			time.Sleep(dispatcherMonitoringPeriod)
		}
	}()

	// If you want to enqueue from other places, e.g. in other function,
	// just pass Dispatcher as Enqueuer
	// (notice that sendJobs() receives jodi.Enqueuer).
	sendJobs(dp)

	// Just helper code to not to exit early.
	stop := make(chan os.Signal, 8)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
}

func sendJobs(enq jodi.Enqueuer) {
	// Now you can add any Job (see printerJob below as an example implementation)
	// to a queue, and it will not block even if there is no available worker currently.
	for i := 0; i < numberOfJobs; i++ {
		log.Printf("[SENDER ] Sending job %d ...\n", i)
		time.Sleep(senderPeriod)

		enq.Enqueue(printerJob{
			number:   i,
			duration: time.Millisecond * 100 * time.Duration(rand.Intn(30)),
		})

		log.Println("[SENDER ] Job sent!")
	}
}

type printerJob struct {
	number   int
	duration time.Duration
}

func (j printerJob) Run() error {
	log.Printf("[JOB #%2d] Doing job for %s...\n", j.number, j.duration.String())
	time.Sleep(j.duration)
	log.Printf("[JOB #%2d] Done job!\n", j.number)
	return nil
}
