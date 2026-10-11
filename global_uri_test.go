package mongodb_locks

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlobalUriLockFactory_InitiallyInitialized(t *testing.T) {
	assert.NotNil(t, uriStorageLockFactoryBeanFactory)
}

func TestNewMongoLockByUri_InvalidUri(t *testing.T) {
	_, err := NewMongoLockByUri(context.Background(), "invalid-uri", "test-lock")
	assert.Error(t, err)
}

func TestGetMongoLockFactoryByUri_InvalidUri(t *testing.T) {
	_, err := GetMongoLockFactoryByUri(context.Background(), "invalid-uri")
	assert.Error(t, err)
}

func TestNewMongoLockByUriWithOptions_InvalidUri(t *testing.T) {
	_, err := NewMongoLockByUriWithOptions(context.Background(), "invalid-uri", nil)
	assert.Error(t, err)
}

