package core

import (
	"context"
	"sync"
	"time"

	"github.com/ivan1993spb/availabilitychecker/internal/log"
	"github.com/ivan1993spb/availabilitychecker/internal/model"
)

type Fetcher interface {
	Tasks(ctx context.Context) <-chan *model.Task
}

type RandomFetcher struct{}

func (RandomFetcher) Tasks(ctx context.Context) <-chan *model.Task {
	const size = 30

	ch := make(chan *model.Task, size)
	defer close(ch)

	for i := 0; i < size; i++ {
		ch <- model.NewTaskRandom()
	}

	return ch
}

type Generator struct {
	run sync.Once

	out   chan *model.Task
	fetch Fetcher
}

func NewGenerator(f Fetcher) *Generator {
	const outBufferSize = 1024

	return &Generator{
		out:   make(chan *model.Task, outBufferSize),
		fetch: f,
	}
}

func (g *Generator) Out() <-chan *model.Task {
	return g.out
}

func (g *Generator) generate(ctx context.Context) {
	ch := g.fetch.Tasks(ctx)

	for {
		select {
		case task, ok := <-ch:
			if !ok {
				return
			}

			select {
			case g.out <- task:

			case <-ctx.Done():
				return
			}

		case <-ctx.Done():
			return
		}
	}
}

func (g *Generator) Run(ctx context.Context) {
	const delay = time.Second * 30

	g.run.Do(func() {
		defer close(g.out)

		log.FromContext(ctx).WithField("delay", delay).
			Debug("generator started")

		ticker := time.NewTicker(delay)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				g.generate(ctx)

			case <-ctx.Done():
				return
			}
		}
	})
}
