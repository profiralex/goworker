package goworker

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

// Worker performs periodic task at certain intervals
type Worker struct {
	name      string
	interval  time.Duration
	action    func() error
	stopChan  chan bool
	errorChan chan error
}

// CreateWorker creates a Worker to run at certain interval and to perform a certain action
func CreateWorker(name string, interval time.Duration, action func() error) *Worker {
	worker := Worker{}
	worker.interval = interval
	worker.action = action
	worker.name = name

	return &worker
}

// Start the Worker
func (worker *Worker) Start() (<-chan error, error) {
	if worker.stopChan != nil {
		return nil, fmt.Errorf("Worker already started")
	}

	worker.stopChan = make(chan bool)
	worker.errorChan = make(chan error)

	defer worker.runLoop()

	return worker.errorChan, nil
}

// Stop the Worker
func (worker *Worker) Stop() {
	if worker.stopChan == nil {
		return
	}

	worker.stopChan <- true
}

func (worker *Worker) runLoop() {
	go worker.run()
}

func (worker *Worker) run() {
	err := worker.performAction()
	if err == nil {
		err = worker.loop()
	}

	worker.errorChan <- err
	if err != nil {
		log.Errorf("Worker failed: %s", err)
	}
	worker.closeChans()
}

func (worker *Worker) loop() error {
	var err error
	timer := time.NewTicker(worker.interval)
	run := true
	for run {
		select {
		case <-worker.stopChan:
			run = false

		case <-timer.C:
			err = worker.performAction()
			if err != nil {
				err = fmt.Errorf("Worker loop failed: %w", err)
				run = false
			}

		default:
			// acceptable delay in case of stop opration
			time.Sleep(50 * time.Millisecond)
		}

	}
	timer.Stop()
	return err
}

func (worker *Worker) performAction() error {
	err := worker.action()
	if err != nil {
		return fmt.Errorf("Worker action failed: %w", err)
	}
	return nil
}

func (worker *Worker) closeChans() {
	close(worker.stopChan)
	worker.stopChan = nil
	close(worker.errorChan)
	worker.errorChan = nil
}
