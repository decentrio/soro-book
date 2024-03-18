package service

import (
	"errors"
)

var (
	// ErrAlreadyStarted is returned when somebody tries to start an already
	// running service.
	ErrAlreadyStarted = errors.New("already started")
	// ErrNotStarted is returned when somebody tries to stop a not running
	// service.
	ErrNotStarted = errors.New("not started")
)

type Service interface {
	// Start services
	Start() error
	OnStart() error

	// Stop services
	Stop() error
	OnStop() error

	// Check if services is running
	IsRunning() bool

	// Terminate services
	Terminate() <-chan struct{}
}

type BaseService struct {
	name      string
	isStarted bool
	terminate chan struct{}

	impl Service
}

// NewBaseService creates a new BaseService.
func NewBaseService(name string, impl Service) *BaseService {
	return &BaseService{
		name:      name,
		terminate: make(chan struct{}),
		isStarted: false,
		impl:      impl,
	}
}

// Start servies
func (bs *BaseService) Start() error {
	if bs.isStarted {
		return ErrAlreadyStarted
	}
	bs.isStarted = true
	return bs.impl.OnStart()
}

// Stop services
func (bs *BaseService) Stop() error {
	if !bs.isStarted {
		return ErrNotStarted
	}
	bs.isStarted = false
	close(bs.terminate)
	return bs.impl.OnStop()
}

// IsRunning() check if services is running
func (bs *BaseService) IsRunning() bool {
	return bs.isStarted
}

// Terminate() return chan if services is terminated
func (bs *BaseService) Terminate() <-chan struct{} {
	return bs.terminate
}
