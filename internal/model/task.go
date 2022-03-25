package model

import (
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/icrowley/fake"
)

type Task struct {
	CheckID uint

	// Host to be checked
	Host string

	// Port to be checked
	Port uint16

	// Timeout
	Timeout time.Duration

	// Number of attempts
	Attempts int

	// Number of successful attempts
	Successes int

	// Last occurred error or nil
	LastError error
}

func NewTaskRandom() *Task {
	return &Task{
		CheckID: uint(rand.Intn(10000)),
		Host:    fake.IPv4(),
		Port:    uint16(rand.Intn(10000)),
		Timeout: 1000,
	}
}

func (t *Task) Addr() string {
	port := strconv.Itoa(int(t.Port))

	return net.JoinHostPort(t.Host, port)
}
