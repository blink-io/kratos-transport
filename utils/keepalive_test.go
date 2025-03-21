package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeepAliveService(t *testing.T) {
	svc := NewKeepAliveService(nil)
	assert.NotNil(t, svc)
	err := svc.Start()
	assert.Nil(t, err)
}
