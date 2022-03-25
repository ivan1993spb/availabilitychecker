package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/ivan1993spb/availabilitychecker/internal/log"
)

func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: log.NewGormLogger(),
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}
