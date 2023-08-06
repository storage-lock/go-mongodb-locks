package mariadb_locks

import (
	"context"
	mongodb_storage "github.com/storage-lock/go-mongodb-storage"
	storage_lock_factory "github.com/storage-lock/go-storage-lock-factory"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoLockFactory struct {
	*storage_lock_factory.StorageLockFactory[*mongo.Client]
}

func NewMongoLockFactory(uri string) (*MongoLockFactory, error) {

	connectionManager := mongodb_storage.NewMongoConnectionManager(uri)
	options := mongodb_storage.NewMongoStorageOptions().SetConnectionManager(connectionManager)
	mongoStorage, err := mongodb_storage.NewMongoStorage(context.Background(), options)
	if err != nil {
		return nil, err
	}
	factory := storage_lock_factory.NewStorageLockFactory[*mongo.Client](mongoStorage, nil)
	return &MongoLockFactory{
		StorageLockFactory: factory,
	}, nil
}
