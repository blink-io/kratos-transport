package utils

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestKeepAliveService(t *testing.T) {
	svc := NewKeepAliveService(nil)
	assert.NotNil(t, svc)

	var wg sync.WaitGroup
	wg.Add(1)

	var startErr error
	go func() {
		startErr = svc.Start()
		wg.Done()
	}()

	// Wait for listener to be ready by checking endpoint
	for i := 0; i < 50; i++ {
		svc.mu.Lock()
		lis := svc.lis
		svc.mu.Unlock()
		if lis != nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	// Verify listener was created
	svc.mu.Lock()
	assert.NotNil(t, svc.lis)
	lis := svc.lis
	svc.mu.Unlock()
	assert.Equal(t, "tcp", lis.Addr().Network())

	// Stop the server
	err := svc.Stop(context.Background())
	assert.Nil(t, err)

	wg.Wait()
	assert.Nil(t, startErr)
}

func TestKeepAliveServicePortRange(t *testing.T) {
	svc := NewKeepAliveService(nil)
	err := svc.generateEndpoint()
	assert.Nil(t, err)

	svc.mu.Lock()
	defer svc.mu.Unlock()
	assert.NotNil(t, svc.lis)

	addr := svc.lis.Addr().(*net.TCPAddr)
	assert.GreaterOrEqual(t, addr.Port, 10000)
	assert.Less(t, addr.Port, 65535)

	svc.lis.Close()
	svc.lis = nil
	svc.endpoint = nil
}
