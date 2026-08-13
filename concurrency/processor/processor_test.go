package processor

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type mockTask struct {
	doFunc func()
}

func (m *mockTask) Do() {
	if m.doFunc != nil {
		m.doFunc()
	}
}

func TestProcessor_Success(t *testing.T) {
	var count int32
	p := NewProcessor(2, 5)

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		task := &mockTask{
			doFunc: func() {
				atomic.AddInt32(&count, 1)
				wg.Done()
			},
		}
		err := p.Do(task)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	}

	wg.Wait()
	p.Close()

	if atomic.LoadInt32(&count) != 5 {
		t.Errorf("expected 5 tasks to be completed, got %d", count)
	}
}

func TestProcessor_QueueFull(t *testing.T) {
	// 1 worker, 1 queue size.
	// We'll submit a blocking task to occupy the worker.
	// Then a task to occupy the queue.
	// The third task should fail with ErrQueueFull.
	p := NewProcessor(1, 1)

	blocked := make(chan struct{})
	doneBlocked := make(chan struct{})

	task1 := &mockTask{
		doFunc: func() {
			<-blocked
			close(doneBlocked)
		},
	}

	// This task starts execution immediately.
	if err := p.Do(task1); err != nil {
		t.Fatalf("first task Do failed: %v", err)
	}

	// Wait a moment for worker to take task1.
	time.Sleep(10 * time.Millisecond)

	task2 := &mockTask{}
	// This task goes into the queue.
	if err := p.Do(task2); err != nil {
		t.Fatalf("second task Do failed: %v", err)
	}

	task3 := &mockTask{}
	// The queue is full, so this should return ErrQueueFull.
	if err := p.Do(task3); err != ErrQueueFull {
		t.Errorf("expected ErrQueueFull, got %v", err)
	}

	// Unblock first task and close.
	close(blocked)
	<-doneBlocked
	p.Close()
}

func TestProcessor_Close(t *testing.T) {
	p := NewProcessor(2, 5)

	task := &mockTask{}
	p.Close()

	// Should return ErrClosed after close
	if err := p.Do(task); err != ErrClosed {
		t.Errorf("expected ErrClosed, got %v", err)
	}
}

func TestProcessor_Race(t *testing.T) {
	p := NewProcessor(4, 100)
	var wg sync.WaitGroup

	// Run multiple goroutines submitting tasks and closing the processor concurrently
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				task := &mockTask{}
				err := p.Do(task)
				if err == ErrClosed {
					return
				}
				time.Sleep(time.Microsecond)
			}
		}()
	}

	time.Sleep(50 * time.Millisecond)
	p.Close()
	wg.Wait()
}
