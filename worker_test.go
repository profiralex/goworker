package goworker_test

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"testing"
	"time"

	"github.com/profiralex/goworker"
)

var defaultTimeout = 5 * time.Second

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestWorkerStartStopSuccess(t *testing.T) {
	timeoutWaiter := time.After(defaultTimeout)
	done := make(chan bool)

	go func() {
		var worker *goworker.Worker
		testAction := func() error {
			return nil
		}
		worker = goworker.CreateWorker(150*time.Millisecond, testAction)

		errChan, err := worker.Start()
		if err != nil {
			t.Errorf("Failed to start worker: %s", err)
		}

		// Stop worker after 1 second
		go func() {
			time.Sleep(1 * time.Second)
			worker.Stop()
		}()

		err = <-errChan
		if err != nil {
			t.Errorf("Worker unexpectedly failed: %s", err)
		}

		done <- true
	}()

	select {
	case <-timeoutWaiter:
		t.Fatal("TIMEOUT!!! Test exceeded expected execution time")
	case <-done:
	}
}

func TestWorkerStartStopError(t *testing.T) {
	timeoutWaiter := time.After(defaultTimeout)
	done := make(chan bool)

	go func() {
		var worker *goworker.Worker
		testAction := func() error {
			return fmt.Errorf("random error")
		}
		worker = goworker.CreateWorker(150*time.Millisecond, testAction)

		errChan, err := worker.Start()
		if err != nil {
			t.Errorf("Failed to start worker: %s", err)
		}

		// Stop worker after 1 second
		go func() {
			time.Sleep(1 * time.Second)
			worker.Stop()
		}()

		err = <-errChan
		if err == nil {
			t.Errorf("Expected to receive an error")
		}

		done <- true
	}()

	select {
	case <-timeoutWaiter:
		t.Fatal("TIMEOUT!!! Test exceeded expected execution time")
	case <-done:
	}
}
