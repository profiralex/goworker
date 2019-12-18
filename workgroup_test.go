package goworker_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/profiralex/goworker"
)

func TestWorkGroupStartStopSuccess(t *testing.T) {
	timeoutWaiter := time.After(defaultTimeout)
	done := make(chan bool)

	go func() {
		testAction := func() error {
			return nil
		}
		group := goworker.CreateWorkGroup(
			goworker.CreateWorker("test worker 1", 150*time.Millisecond, testAction),
			goworker.CreateWorker("test worker 2", 150*time.Millisecond, testAction),
		)

		doneChan, err := group.Start()
		if err != nil {
			t.Errorf("Failed to start worker: %s", err)
		}

		// Stop worker after 1 second
		go func() {
			time.Sleep(1 * time.Second)
			group.Stop()
		}()

		_ = <-doneChan

		done <- true
	}()

	select {
	case <-timeoutWaiter:
		t.Fatal("TIMEOUT!!! Test exceeded expected execution time")
	case <-done:
	}
}

func TestWorkGroupStartStopError(t *testing.T) {
	timeoutWaiter := time.After(defaultTimeout)
	done := make(chan bool)

	go func() {
		testAction := func() error {
			return fmt.Errorf("random error")
		}
		group := goworker.CreateWorkGroup(
			goworker.CreateWorker("test worker 1", 150*time.Millisecond, testAction),
			goworker.CreateWorker("test worker 2", 150*time.Millisecond, testAction),
		)

		doneChan, err := group.Start()
		if err != nil {
			t.Errorf("Failed to start worker: %s", err)
		}

		// Stop worker after 1 second
		go func() {
			time.Sleep(1 * time.Second)
			group.Stop()
		}()

		_ = <-doneChan

		done <- true
	}()

	select {
	case <-timeoutWaiter:
		t.Fatal("TIMEOUT!!! Test exceeded expected execution time")
	case <-done:
	}
}
