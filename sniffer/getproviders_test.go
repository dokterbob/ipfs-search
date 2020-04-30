package sniffer

import (
	"context"
	"errors"
	"github.com/ipfs-search/ipfs-search/types"
	"github.com/stretchr/testify/assert"
	"testing"

	"time"
)

const longTime = time.Duration(time.Minute)
const shortTime = time.Duration(0)

type mockLogger struct {
	wait time.Duration
	msg  map[string]interface{}
	err  error
}

func (m mockLogger) Next() (map[string]interface{}, error) {
	time.Sleep(m.wait)

	return m.msg, m.err
}

type mockProviderExtractor struct {
	provider *types.Provider
	err      error
}

func (m mockProviderExtractor) Extract(map[string]interface{}) (*types.Provider, error) {
	return m.provider, m.err
}

// TestContextCancel tests whether we're returning an error on context cancellation
func TestContextCancel(t *testing.T) {
	l := &mockLogger{
		wait: time.Duration(longTime),
	}

	e := &mockProviderExtractor{}

	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan types.Provider)

	assert := assert.New(t)

	// Cancel context immediately
	cancel()

	err := getProviders(ctx, l, e, c, longTime)

	assert.Equal(err, context.Canceled)
}

// TestTimeout tests whether we're returning an error on timeout
func TestTimeout(t *testing.T) {
	l := &mockLogger{
		wait: time.Duration(longTime),
	}

	e := &mockProviderExtractor{}

	ctx := context.Background()
	c := make(chan types.Provider)

	assert := assert.New(t)

	err := getProviders(ctx, l, e, c, shortTime)

	assert.Equal(ErrorLoggerTimeout, err)

}

// TestYieldProvider tests whether a provider is yielded on a provider message
func TestYieldProvider(t *testing.T) {
	mockProvider := types.Provider{
		Resource: &types.Resource{
			Protocol: "mockProtocol",
			Id:       "mockKey",
		},
	}

	l := &mockLogger{}

	e := &mockProviderExtractor{
		provider: &mockProvider,
	}

	ctx := context.Background()
	c := make(chan types.Provider)

	assert := assert.New(t)

	go getProviders(ctx, l, e, c, longTime)

	provider := <-c
	assert.Equal(mockProvider, provider)
}

// TestLoggerError tests for error propagation from the logger
func TestLoggerError(t *testing.T) {
	errMock := errors.New("mock")

	l := &mockLogger{
		err: errMock,
	}

	e := &mockProviderExtractor{}

	ctx := context.Background()
	c := make(chan types.Provider)

	assert := assert.New(t)

	err := getProviders(ctx, l, e, c, longTime)

	assert.Equal(err, errMock)
}

// TestProviderErr tests for error propagation from ResourceProvider
func TestProviderErr(t *testing.T) {
	errMock := errors.New("mock")

	l := &mockLogger{}

	e := &mockProviderExtractor{
		err: errMock,
	}

	ctx := context.Background()
	c := make(chan types.Provider)

	assert := assert.New(t)

	err := getProviders(ctx, l, e, c, longTime)

	assert.Equal(err, errMock)
}
