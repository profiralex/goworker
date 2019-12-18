package goworker

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"sync"
)

// WorkGroup helps to handle multiple workers
type WorkGroup struct {
	workers  []*Worker
	doneChan chan bool
	started  bool
}

// CreateWorkGroup creates a WorkGroup to run workers
func CreateWorkGroup(workers ...*Worker) *WorkGroup {
	return &WorkGroup{
		workers: workers,
	}
}

// Start the WorkGroup
func (group *WorkGroup) Start() (<-chan bool, error) {
	if group.doneChan != nil {
		return nil, fmt.Errorf("WorkGroup already started")
	}

	group.doneChan = make(chan bool)

	defer group.run()

	return group.doneChan, nil
}

// Stop the WorkGroup
func (group *WorkGroup) Stop() {
	if group.doneChan == nil {
		return
	}

	for _, worker := range group.workers {
		worker.Stop()
	}
}

func (group *WorkGroup) run() {
	go group.runWorkers()
}

func (group *WorkGroup) runWorkers() {
	wg := new(sync.WaitGroup)
	for _, worker := range group.workers {
		wg.Add(1)
		go group.runWorker(wg, worker)
	}

	wg.Wait()
	group.doneChan <- true
	group.doneChan = nil
}

func (group *WorkGroup) runWorker(wg *sync.WaitGroup, worker *Worker) {
	defer wg.Done()
	errChan, err := worker.Start()
	if err != nil {
		log.Errorf("Failed to start worker %s : %s", worker.name, err)
		return
	}

	log.Infof("Started worker %s", worker.name)

	err = <-errChan
	if err != nil {
		log.Errorf("Worker %s failed : %s", worker.name, err)
		return
	}

	log.Infof("Worker %s stopped", worker.name)
}
