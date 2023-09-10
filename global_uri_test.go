package mariadb_locks

import (
	"context"
	storage_lock_test_helper "github.com/storage-lock/go-storage-lock-test-helper"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const MONGO_URI_ENV_NAME = "STORAGE_LOCK_MONGO_URI"

func TestNewMongoLockByUri(t *testing.T) {

	mongoUri := os.Getenv(MONGO_URI_ENV_NAME)
	assert.NotEmpty(t, mongoUri)

	factory, err := GetMongoLockFactoryByUri(context.Background(), mongoUri)
	assert.Nil(t, err)

	storage_lock_test_helper.PlayerNum = 10
	storage_lock_test_helper.EveryOnePlayTimes = 100
	storage_lock_test_helper.TestStorageLock(t, factory)
}
