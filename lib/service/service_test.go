package service_test

import (
	"testing"

	"github.com/decentrio/soro-book/lib/service"
	"github.com/stretchr/testify/require"
)

type testService struct {
	service.BaseService
}

func (ms testService) OnStart() error {
	return nil
}

func (ms testService) OnStop() error {
	return nil
}
func TestBaseServiceStart(t *testing.T) {
	ts := &testService{}
	ts.BaseService = *service.NewBaseService("TestService", ts)
	err := ts.Start()
	require.NoError(t, err)
	require.True(t, ts.BaseService.IsRunning())
}

func TestBaseServiceStop(t *testing.T) {
	ts := &testService{}
	ts.BaseService = *service.NewBaseService("TestService", ts)
	err := ts.Start()
	require.NoError(t, err)
	require.True(t, ts.BaseService.IsRunning())

	err = ts.Stop()
	require.NoError(t, err)
	require.False(t, ts.BaseService.IsRunning())
}
