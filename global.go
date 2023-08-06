package mariadb_locks

import (
	storage_lock "github.com/storage-lock/go-storage-lock"
	"sync"
)

var GlobalMongoLockFactory *MongoLockFactory
var globalMongoLockFactoryOnce sync.Once
var globalMongoLockFactoryErr error

func InitGlobalMongoLockFactory(uri string) error {
	factory, err := NewMongoLockFactory(uri)
	if err != nil {
		return err
	}
	GlobalMongoLockFactory = factory
	return nil
}

func NewMongoLock(uri string, lockId string) (*storage_lock.StorageLock, error) {
	globalMongoLockFactoryOnce.Do(func() {
		globalMongoLockFactoryErr = InitGlobalMongoLockFactory(uri)
	})
	if globalMongoLockFactoryErr != nil {
		return nil, globalMongoLockFactoryErr
	}
	return GlobalMongoLockFactory.CreateLock(lockId)
}

//func NewMongoLockFromClient(client *mongo.Client, lockId string) (*storage_lock.StorageLock, error) {
//	globalMongoLockFactoryOnce.Do(func() {
//		globalMongoLockFactoryErr = InitGlobalMongoLockFactory(uri)
//	})
//	if globalMongoLockFactoryErr != nil {
//		return nil, globalMongoLockFactoryErr
//	}
//	return GlobalMongoLockFactory.CreateLock(lockId)
//}
