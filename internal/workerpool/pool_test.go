package workerpool

import (
	"context"
	"testing"

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
