package workerpool

import (
	"context"
	"sync"
)

type Job struct {
	ID      int
	Content []byte
	Func    func([]byte) ([]byte, error)
}

type Result struct {
	JobID   int
	Content []byte
	Error   error
}

type Pool struct {
	workerCount int
	jobs        chan Job
	results     chan Result
	wg          sync.WaitGroup
}

func NewPool(wc int, buffer int) *Pool {
	return &Pool{
		workerCount: wc,
		jobs:        make(chan Job, buffer),
		results:     make(chan Result, buffer),
	}
}

func (p *Pool) worker(ctx context.Context) {
	defer p.wg.Done()
	for {
		select {
		case job, ok := <-p.jobs:
			if !ok {
				return
			}
			result, err := job.Func(job.Content)
			if err != nil {
				p.results <- Result{
					JobID: job.ID,
					Error: err,
				}
				continue
			}
			p.results <- Result{
				JobID:   job.ID,
				Content: result,
			}

		case <-ctx.Done():
			return
		}

	}

}

func (p *Pool) Start(ctx context.Context) <-chan Result {
	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)
		go p.worker(ctx)
	}

	return p.results
}

func (p *Pool) Submit(job Job) {
	p.jobs <- job
}

func (p *Pool) Shutdown() {
	close(p.jobs)
	p.wg.Wait()
	close(p.results)
}
