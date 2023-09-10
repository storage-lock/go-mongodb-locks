package mongodb_locks

import (
	"context"
	mongodb_storage "github.com/storage-lock/go-mongodb-storage"
	storage_lock "github.com/storage-lock/go-storage-lock"
	storage_lock_factory "github.com/storage-lock/go-storage-lock-factory"
	"go.mongodb.org/mongo-driver/mongo"
)

var clientStorageLockFactoryBeanFactory *storage_lock_factory.StorageLockFactoryBeanFactory[*mongo.Client, *mongo.Client] = storage_lock_factory.NewStorageLockFactoryBeanFactory[*mongo.Client, *mongo.Client]()

func NewMongoLockByClient(ctx context.Context, client *mongo.Client, lockId string) (*storage_lock.StorageLock, error) {
	factory, err := GetMongoLockFactoryByClient(ctx, client)
	if err != nil {
		return nil, err
	}
	return factory.CreateLock(lockId)
}

func NewMongoLockByClientWithOptions(ctx context.Context, client *mongo.Client, options *storage_lock.StorageLockOptions) (*storage_lock.StorageLock, error) {
	factory, err := GetMongoLockFactoryByClient(ctx, client)
	if err != nil {
		return nil, err
	}
	return factory.CreateLockWithOptions(options)
}

func GetMongoLockFactoryByClient(ctx context.Context, client *mongo.Client) (*storage_lock_factory.StorageLockFactory[*mongo.Client], error) {
	return clientStorageLockFactoryBeanFactory.GetOrInit(ctx, client, func(ctx context.Context) (*storage_lock_factory.StorageLockFactory[*mongo.Client], error) {
		connectionManager := mongodb_storage.NewMongoConnectionManagerFromClient(client)
		options := mongodb_storage.NewMongoStorageOptions().SetConnectionManager(connectionManager)
		storage, err := mongodb_storage.NewMongoStorage(ctx, options)
		if err != nil {
			return nil, err
		}
		factory := storage_lock_factory.NewStorageLockFactory(storage, options.ConnectionManager)
		return factory, nil
	})
}
