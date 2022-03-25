package core

import (
	"context"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/icrowley/fake"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/availabilitychecker/internal/log"
	"github.com/ivan1993spb/availabilitychecker/internal/model"
)

type Executer interface {
	Execute(ctx context.Context, task *model.Task) error
}

type RandExecutor struct{}

func (e RandExecutor) Execute(ctx context.Context, task *model.Task) error {
	success := rand.Intn(2) == 0

	if !success {
		color := fake.Color()
		err := errors.Errorf("%s error", color)

		fields := logrus.Fields{
			logrus.ErrorKey: err,
			"address":       task.Addr(),
		}

		log.FromContext(ctx).WithFields(fields).
			Debug("failed to connect")

		return err
	}

	return nil
}

type CheckExecutor struct {
}

func (CheckExecutor) Execute(ctx context.Context, task *model.Task) error {
	addr := task.Addr()

	var dialer net.Dialer
	_, err := dialer.DialContext(ctx, "tcp", addr)

	if err != nil {
		fields := logrus.Fields{
			logrus.ErrorKey: err,
			"address":       addr,
		}

		log.FromContext(ctx).WithFields(fields).
			Debug("failed to connect")

		return err
	}

	return nil
}

type Pool struct {
	run sync.Once
	wg  sync.WaitGroup

	in    <-chan *model.Task
	round chan *model.Task
	out   chan *model.Result

	exec Executer
}

func NewPool(exec Executer, in <-chan *model.Task) *Pool {
	const (
		outBufferSize   = 128
		roundBufferSize = 512
	)

	return &Pool{
		in:    in,
		out:   make(chan *model.Result, outBufferSize),
		round: make(chan *model.Task, roundBufferSize),

		exec: exec,
	}
}

const (
	AttemptLimit      = 3
	SuccessThreashold = 2
)

const msgEmptyFailMessage = "empty fail message"

func (p *Pool) execute(ctx context.Context, task *model.Task) {
	const sendTimeout = time.Millisecond * 100

	task.Attempts += 1

	ctxExec, cancel := context.WithTimeout(ctx, task.Timeout)
	defer cancel()

	err := p.exec.Execute(ctxExec, task)

	if err != nil {
		task.LastError = err
	} else {
		task.Successes += 1
	}

	if task.Attempts < AttemptLimit {
		select {
		case p.round <- task:

		case <-ctx.Done():
		}

		return
	}

	result := &model.Result{
		Task: task,
	}

	if task.Successes < SuccessThreashold {
		if task.LastError != nil {
			result.FailMessage = task.LastError.Error()
		} else {
			result.FailMessage = msgEmptyFailMessage
		}
		result.Status = model.StatusFail
	} else {
		result.Status = model.StatusOK
		result.FailMessage = ""
	}

	ctxSend, cancel := context.WithTimeout(ctx, sendTimeout)
	defer cancel()

	select {
	case p.out <- result:

	case <-ctxSend.Done():
	}
}

func (p *Pool) startWorker(ctx context.Context) {
	for {
		select {
		case task := <-p.in:
			p.execute(ctx, task)

		case task := <-p.round:
			p.execute(ctx, task)

		case <-ctx.Done():
			return
		}
	}
}

func (p *Pool) Out() <-chan *model.Result {
	return p.out
}

const minWorkersNumber = 2

func (p *Pool) Run(ctx context.Context, workers int) {
	if workers < minWorkersNumber {
		panic("not enough workers")
	}

	p.run.Do(func() {
		log.FromContext(ctx).WithField("workers", workers).
			Debug("pool started")

		p.wg.Add(workers)

		for i := 0; i < workers; i++ {
			go func() {
				defer p.wg.Done()

				p.startWorker(ctx)
			}()
		}

		p.wg.Wait()

		close(p.round)
		close(p.out)
	})
}
