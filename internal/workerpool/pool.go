package workerpool

import (
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

func NewPool(wc int) *Pool {
	return &Pool{
		workerCount: wc,
		jobs:        make(chan Job, wc),
		results:     make(chan Result, wc),
	}
}

func (p *Pool) worker() {
	defer p.wg.Done()
	for job := range p.jobs {
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
	}
}

func (p *Pool) Start() {
	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)
		go p.worker()
	}
}

func (p *Pool) Submit(job Job) {
	p.jobs <- job
}

func (p *Pool) Shutdown() {
	close(p.jobs)
	p.wg.Wait()
	close(p.results)
}
