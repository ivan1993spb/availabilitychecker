package mysql

import (
	"context"

	"gorm.io/gorm"

	"github.com/ivan1993spb/availabilitychecker/internal/model"
)

type Store struct {
	DB *gorm.DB
}

func (s *Store) Save(ctx context.Context, result *model.Result) {
	check := convertResultToCheck(result)

	s.DB.WithContext(ctx).Model(check).Updates(Check{
		Status:      check.Status,
		FailMessage: check.FailMessage,
	})
}
