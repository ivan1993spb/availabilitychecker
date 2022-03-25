package app

import (
	"context"
	"sync"

	"github.com/ivan1993spb/availabilitychecker/internal/core"
	"github.com/ivan1993spb/availabilitychecker/internal/log"
)

type Config struct {
	Fetcher  core.Fetcher
	Executer core.Executer
	Store    core.Store
}

type App struct {
	run sync.Once
	wg  sync.WaitGroup

	cfg Config
}

func NewApp(cfg Config) *App {
	return &App{
		cfg: cfg,
	}
}

func (a *App) Run(ctx context.Context) {
	a.run.Do(func() {
		log.FromContext(ctx).Info("start")

		generator := core.NewGenerator(a.cfg.Fetcher)
		a.wg.Add(1)
		go func() {
			defer a.wg.Done()
			generator.Run(ctx)
		}()

		pool := core.NewPool(a.cfg.Executer, generator.Out())
		a.wg.Add(1)
		go func() {
			defer a.wg.Done()
			pool.Run(ctx, 10)
		}()

		acceptor := core.NewAcceptor(a.cfg.Store, pool.Out())
		a.wg.Add(1)
		go func() {
			defer a.wg.Done()
			acceptor.Run(ctx)
		}()

		a.wg.Wait()

		log.FromContext(ctx).Info("finish")
	})
}
