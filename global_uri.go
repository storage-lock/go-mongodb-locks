package mongodb_locks

import (
	"context"
	mongodb_storage "github.com/storage-lock/go-mongodb-storage"
	storage_lock "github.com/storage-lock/go-storage-lock"
	storage_lock_factory "github.com/storage-lock/go-storage-lock-factory"
	"go.mongodb.org/mongo-driver/mongo"
)

var uriStorageLockFactoryBeanFactory *storage_lock_factory.StorageLockFactoryBeanFactory[string, *mongo.Client] = storage_lock_factory.NewStorageLockFactoryBeanFactory[string, *mongo.Client]()

// newMongoStorageByURI 指向 mongodb_storage.NewMongoStorage，包级变量以便单测注入构造失败，
// 覆盖 GetOrInit initFunc 的错误分支。
var newMongoStorageByURI = mongodb_storage.NewMongoStorage

func NewMongoLockByUri(ctx context.Context, uri string, lockId string) (*storage_lock.StorageLock, error) {
	factory, err := GetMongoLockFactoryByUri(ctx, uri)
	if err != nil {
		return nil, err
	}
	return factory.CreateLock(lockId)
}

func NewMongoLockByUriWithOptions(ctx context.Context, uri string, options *storage_lock.StorageLockOptions) (*storage_lock.StorageLock, error) {
	factory, err := GetMongoLockFactoryByUri(ctx, uri)
	if err != nil {
		return nil, err
	}
	return factory.CreateLockWithOptions(options)
}

func GetMongoLockFactoryByUri(ctx context.Context, uri string) (*storage_lock_factory.StorageLockFactory[*mongo.Client], error) {
	return uriStorageLockFactoryBeanFactory.GetOrInit(ctx, uri, func(ctx context.Context) (*storage_lock_factory.StorageLockFactory[*mongo.Client], error) {
		options := mongodb_storage.NewMongoStorageOptionsWithURI(uri)
		storage, err := newMongoStorageByURI(ctx, options)
		if err != nil {
			return nil, err
		}
		factory := storage_lock_factory.NewStorageLockFactory(storage, options.ConnectionManager)
		return factory, nil
	})
}
