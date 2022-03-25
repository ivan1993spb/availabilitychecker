package mysql

import (
	"math/rand"

	"github.com/icrowley/fake"
)

type Check struct {
	ID uint `gorm:"primarykey"`

	Host string `gorm:"type:varchar(255); not null"`

	Port uint16

	Status uint8 `gorm:"default:0"`

	// Timeout in milliseconds
	Timeout int64

	FailMessage *string `gorm:"type:varchar(255); default:null"`
}

func NewCheckRandom() *Check {
	title := fake.Title()

	return &Check{
		Host:        fake.IPv4(),
		Port:        uint16(rand.Intn(10000)),
		Status:      uint8(rand.Intn(3)),
		Timeout:     1000,
		FailMessage: &title,
	}
}
