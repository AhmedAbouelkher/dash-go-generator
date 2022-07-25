package main

import "math"

type ConsumerFunc[T any, A any] func(e T) A

type Pool[T any, A any] struct {
	Workers  int
	Consumer ConsumerFunc[T, A]
	Data     []T
}

func RunWorkersPool[T any, A any](p *Pool[T, A]) []A {
	batch := len(p.Data)
	workers := p.Workers
	if workers == 0 {
		workers = int(math.Ceil(float64(batch) / 10))
	}
	// In order to use our pool of workers we need to send
	// them work and collect their results. We make 2
	// channels for this.
	jobs := make(chan T, batch)
	results := make(chan A, batch)

	// This starts up n workers, initially blocked
	// because there are no jobs yet.
	for w := 1; w <= workers; w++ {
		go func(jobs <-chan T, results chan<- A, consumer ConsumerFunc[T, A]) {
			for job := range jobs {
				results <- consumer(job)
			}
		}(jobs, results, p.Consumer)
	}

	// Here we send k `jobs` and then `close` that
	// channel to indicate that's all the work we have.
	for _, job := range p.Data {
		jobs <- job
	}
	close(jobs)

	res := make([]A, 0)

	// Finally we collect all the results of the work.
	for i := 0; i < batch; i++ {
		r := <-results
		res = append(res, r)
	}

	return res
}
