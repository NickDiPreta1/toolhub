package workerpool

import (
	"context"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/NickDiPreta1/toolhub/internal/tools/hashutil"
)

func TestPoolSuccess(t *testing.T) {
	ctx := context.Background()
	pool := NewPool(3, 3)
	resChan := pool.Start(ctx)

	data := []byte("Some data")

	job1 := Job{
		ID:      1,
		Content: data,
		Func: func(b []byte) ([]byte, error) {
			result, err := hashutil.Hash(data)
			return []byte(result), err
		},
	}

	pool.Submit(job1)

	var results []Result
	done := make(chan struct{})

	go func() {
		for result := range resChan {
			results = append(results, result)
		}
		close(done)
	}()

	pool.Shutdown()
	<-done

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].JobID != job1.ID {
		t.Errorf("Expected ID %d got %d", job1.ID, results[0].JobID)
	}

	if results[0].Error != nil {
		t.Errorf("Expected no error, got %v", results[0].Error)
	}

	if len(results[0].Content) == 0 {
		t.Error("Expected content in result, got empty")
	}
}

func TestPoolMultipleJobs(t *testing.T) {
	ctx := context.Background()
	pool := NewPool(3, 10)
	resChan := pool.Start(ctx)

	jobCount := 10
	for i := 1; i <= jobCount; i++ {
		data := []byte("data" + string(rune(i)))
		job := Job{
			ID:      i,
			Content: data,
			Func: func(b []byte) ([]byte, error) {
				result, err := hashutil.Hash(b)
				return []byte(result), err
			},
		}
		pool.Submit(job)
	}

	var results []Result
	done := make(chan struct{})

	go func() {
		for result := range resChan {
			results = append(results, result)
		}
		close(done)
	}()

	pool.Shutdown()
	<-done

	if len(results) != jobCount {
		t.Errorf("Expected %d results, got %d", jobCount, len(results))
	}

	for _, result := range results {
		if result.Error != nil {
			t.Errorf("Job %d returned error: %v", result.JobID, result.Error)
		}
		if len(result.Content) == 0 {
			t.Errorf("Job %d returned empty content", result.JobID)
		}
	}
}

func TestPoolWithErrors(t *testing.T) {
	ctx := context.Background()
	pool := NewPool(2, 5)
	resChan := pool.Start(ctx)

	successJob := Job{
		ID:      1,
		Content: []byte("success"),
		Func: func(b []byte) ([]byte, error) {
			result, err := hashutil.Hash(b)
			return []byte(result), err
		},
	}

	errorJob := Job{
		ID:      2,
		Content: []byte("error"),
		Func: func(b []byte) ([]byte, error) {
			return nil, context.DeadlineExceeded
		},
	}

	pool.Submit(successJob)
	pool.Submit(errorJob)

	var results []Result
	done := make(chan struct{})

	go func() {
		for result := range resChan {
			results = append(results, result)
		}
		close(done)
	}()

	pool.Shutdown()
	<-done

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	var successResult, errorResult *Result
	for i := range results {
		if results[i].JobID == 1 {
			successResult = &results[i]
		} else if results[i].JobID == 2 {
			errorResult = &results[i]
		}
	}

	if successResult == nil {
		t.Fatal("Success result not found")
	}
	if successResult.Error != nil {
		t.Errorf("Success job should not have error, got %v", successResult.Error)
	}

	if errorResult == nil {
		t.Fatal("Error result not found")
	}
	if errorResult.Error == nil {
		t.Error("Error job should have error, got nil")
	}
}

func TestPoolContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	pool := NewPool(2, 5)
	resChan := pool.Start(ctx)

	var results []Result
	done := make(chan struct{})

	go func() {
		for result := range resChan {
			results = append(results, result)
		}
		close(done)
	}()

	cancel()

	pool.Shutdown()
	<-done
}

func TestPoolSingleWorker(t *testing.T) {
	ctx := context.Background()
	pool := NewPool(1, 3)
	resChan := pool.Start(ctx)

	for i := 1; i <= 3; i++ {
		job := Job{
			ID:      i,
			Content: []byte("data"),
			Func: func(b []byte) ([]byte, error) {
				result, err := hashutil.Hash(b)
				return []byte(result), err
			},
		}
		pool.Submit(job)
	}

	var results []Result
	done := make(chan struct{})

	go func() {
		for result := range resChan {
			results = append(results, result)
		}
		close(done)
	}()

	pool.Shutdown()
	<-done

	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}
}

func TestPoolLargeBufferedJobs(t *testing.T) {
	ctx := context.Background()
	pool := NewPool(5, 100)
	resChan := pool.Start(ctx)

	jobCount := 100
	for i := 1; i <= jobCount; i++ {
		job := Job{
			ID:      i,
			Content: []byte("large batch"),
			Func: func(b []byte) ([]byte, error) {
				result, err := hashutil.Hash(b)
				return []byte(result), err
			},
		}
		pool.Submit(job)
	}

	var results []Result
	done := make(chan struct{})

	go func() {
		for result := range resChan {
			results = append(results, result)
		}
		close(done)
	}()

	pool.Shutdown()
	<-done

	if len(results) != jobCount {
		t.Errorf("Expected %d results, got %d", jobCount, len(results))
	}

	for _, result := range results {
		if result.Error != nil {
			t.Errorf("Job %d returned error: %v", result.JobID, result.Error)
		}
	}
}

func TestPoolNoGoroutineLeak(t *testing.T) {
	// Force GC to clean up any lingering goroutines from previous tests
	runtime.GC()
	time.Sleep(10 * time.Millisecond)

	// Get initial goroutine count
	initialGoroutines := runtime.NumGoroutine()

	ctx := context.Background()
	workerCount := 5
	pool := NewPool(workerCount, 10)
	resChan := pool.Start(ctx)

	// Submit several jobs
	jobCount := 20
	for i := 1; i <= jobCount; i++ {
		job := Job{
			ID:      i,
			Content: []byte("test data"),
			Func: func(b []byte) ([]byte, error) {
				result, err := hashutil.Hash(b)
				return []byte(result), err
			},
		}
		pool.Submit(job)
	}

	// Collect results
	var results []Result
	done := make(chan struct{})

	go func() {
		for result := range resChan {
			results = append(results, result)
		}
		close(done)
	}()

	// Shutdown pool and wait for completion
	pool.Shutdown()
	<-done

	// Give goroutines time to fully terminate
	runtime.GC()
	time.Sleep(50 * time.Millisecond)

	// Check final goroutine count
	finalGoroutines := runtime.NumGoroutine()

	// The goroutine count should return to the initial count
	// We allow a small delta for test framework variance
	leakedGoroutines := finalGoroutines - initialGoroutines
	if leakedGoroutines > 0 {
		t.Errorf("Goroutine leak detected: started with %d, ended with %d (%d leaked)",
			initialGoroutines, finalGoroutines, leakedGoroutines)
	}

	// Verify all jobs completed successfully
	if len(results) != jobCount {
		t.Errorf("Expected %d results, got %d", jobCount, len(results))
	}
}

func TestPoolConcurrentSubmission(t *testing.T) {
	ctx := context.Background()
	pool := NewPool(10, 200)
	resChan := pool.Start(ctx)

	// Track results with a thread-safe map
	resultCount := make(map[int]bool)
	var resultMu sync.Mutex

	done := make(chan struct{})
	go func() {
		for result := range resChan {
			resultMu.Lock()
			resultCount[result.JobID] = true
			if result.Error != nil {
				t.Errorf("Job %d returned error: %v", result.JobID, result.Error)
			}
			resultMu.Unlock()
		}
		close(done)
	}()

	// Submit jobs concurrently from multiple goroutines
	numSubmitters := 10
	jobsPerSubmitter := 20
	totalJobs := numSubmitters * jobsPerSubmitter

	var submitWg sync.WaitGroup
	submitWg.Add(numSubmitters)

	for submitter := 0; submitter < numSubmitters; submitter++ {
		go func(submitterID int) {
			defer submitWg.Done()
			for i := 0; i < jobsPerSubmitter; i++ {
				jobID := submitterID*jobsPerSubmitter + i
				job := Job{
					ID:      jobID,
					Content: []byte("concurrent test"),
					Func: func(b []byte) ([]byte, error) {
						// Simulate some work
						time.Sleep(1 * time.Millisecond)
						result, err := hashutil.Hash(b)
						return []byte(result), err
					},
				}
				pool.Submit(job)
			}
		}(submitter)
	}

	// Wait for all submissions to complete
	submitWg.Wait()

	// Shutdown and wait for all results
	pool.Shutdown()
	<-done

	// Verify all jobs were processed
	resultMu.Lock()
	defer resultMu.Unlock()

	if len(resultCount) != totalJobs {
		t.Errorf("Expected %d unique results, got %d", totalJobs, len(resultCount))
	}

	// Verify all job IDs were processed
	for i := 0; i < totalJobs; i++ {
		if !resultCount[i] {
			t.Errorf("Job %d was not processed", i)
		}
	}
}
