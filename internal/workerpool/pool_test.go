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

	if results[0].JobID != job1.ID {
		t.Errorf("Expected ID %d got %d", job1.ID, results[0].JobID)
	}
}
