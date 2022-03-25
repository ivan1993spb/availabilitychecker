package core

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/availabilitychecker/internal/log"
	"github.com/ivan1993spb/availabilitychecker/internal/model"
)

type Store interface {
	Save(ctx context.Context, result *model.Result)
}

type Ignore struct{}

func (Ignore) Save(ctx context.Context, result *model.Result) {
	log.FromContext(ctx).WithFields(logrus.Fields{
		"id":       result.Task.CheckID,
		"addr":     result.Task.Addr(),
		"status":   result.Status,
		"fail_msg": result.FailMessage,
	}).Infoln("save result")

	// Ignore
}

type Acceptor struct {
	run sync.Once

	in <-chan *model.Result

	store Store
}

func NewAcceptor(store Store, in <-chan *model.Result) *Acceptor {
	return &Acceptor{
		in: in,

		store: store,
	}
}

func (a *Acceptor) accept(ctx context.Context, result *model.Result) {
	a.store.Save(ctx, result)
}

func (a *Acceptor) Run(ctx context.Context) {
	a.run.Do(func() {
		log.FromContext(ctx).Debug("acceptor started")

		for {
			select {
			case result := <-a.in:
				a.accept(ctx, result)

			case <-ctx.Done():
				return
			}
		}
	})
}
