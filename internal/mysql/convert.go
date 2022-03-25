package mysql

import (
	"time"

	"github.com/ivan1993spb/availabilitychecker/internal/model"
)

func convertChecksToTasks(checks []*Check) []*model.Task {
	tasks := make([]*model.Task, len(checks))

	for i, check := range checks {
		tasks[i] = &model.Task{
			CheckID:   check.ID,
			Host:      check.Host,
			Port:      check.Port,
			Timeout:   time.Millisecond * time.Duration(check.Timeout),
			Attempts:  0,
			Successes: 0,
		}
	}

	return tasks
}

func convertResultToCheck(result *model.Result) *Check {
	return &Check{
		ID:          result.Task.CheckID,
		Host:        result.Task.Host,
		Port:        result.Task.Port,
		Status:      uint8(result.Status),
		Timeout:     result.Task.Timeout.Milliseconds(),
		FailMessage: &result.FailMessage,
	}
}
