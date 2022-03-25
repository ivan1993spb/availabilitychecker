package mysql

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/ivan1993spb/availabilitychecker/internal/log"
	"github.com/ivan1993spb/availabilitychecker/internal/model"
)

type Source struct {
	DB *gorm.DB
}

type TasksBatchFunc func(tasks []*model.Task, batch int) error

const (
	errMsgFetchTasks       = "failed to fetch tasks"
	errMsgHandleTasksBatch = "failed to handle tasks batch"
)

func (s *Source) fetchTasks(ctx context.Context,
	f TasksBatchFunc) error {

	const batchSize = 128

	var checks []*Check

	handle := func(tx *gorm.DB, batch int) error {
		tasks := convertChecksToTasks(checks)

		err := f(tasks, batch)
		if err != nil {
			return errors.Wrap(err, errMsgHandleTasksBatch)
		}

		return nil
	}

	result := s.DB.WithContext(ctx).FindInBatches(&checks, batchSize, handle)
	if result.Error != nil {
		return errors.Wrap(result.Error, errMsgFetchTasks)
	}

	return nil
}

func (s *Source) Tasks(ctx context.Context) <-chan *model.Task {
	const outBuffer = 1024

	out := make(chan *model.Task, outBuffer)

	go func() {
		defer close(out)

		err := s.fetchTasks(ctx, func(tasks []*model.Task, _ int) error {
			for _, task := range tasks {
				select {
				case out <- task:

				case <-ctx.Done():
					return ctx.Err()
				}
			}

			return nil
		})

		if err != nil {
			log.FromContext(ctx).WithError(err).Error("cannot fetch tasks")
			return
		}
	}()

	return out
}
