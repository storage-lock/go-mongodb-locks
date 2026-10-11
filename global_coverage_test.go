package mongodb_locks

import (
	"context"
	"errors"
	"testing"
	"time"

	mongodb_storage "github.com/storage-lock/go-mongodb-storage"
	storage_lock "github.com/storage-lock/go-storage-lock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// newMongoStorage 指向 mongodb_storage.NewMongoStorage 的包级注入点声明在被测包内，
// 此处仅使用；注入变量定义在 global_client.go / global_uri.go 中。

func testMongoURI(t *testing.T) string {
	t.Helper()
	uri := "mongodb://127.0.0.1:27017"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		t.Skipf("mongo 连接失败，跳过: %v", err)
	}
	defer client.Disconnect(ctx)
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		t.Skipf("本机 mongo 不可用，跳过: %v", err)
	}
	return uri
}

func testMongoClient(t *testing.T) *mongo.Client {
	t.Helper()
	uri := testMongoURI(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	assert.Nil(t, err)
	t.Cleanup(func() { client.Disconnect(context.Background()) })
	return client
}

func TestNewMongoLockByClient(t *testing.T) {
	client := testMongoClient(t)
	ctx := context.Background()

	factory, err := GetMongoLockFactoryByClient(ctx, client)
	assert.Nil(t, err)
	assert.NotNil(t, factory)

	// 缓存命中
	factory2, err := GetMongoLockFactoryByClient(ctx, client)
	assert.Nil(t, err)
	assert.Equal(t, factory, factory2)

	lock, err := NewMongoLockByClient(ctx, client, "test-mongo-lock-client")
	assert.Nil(t, err)
	assert.NotNil(t, lock)

	options := storage_lock.NewStorageLockOptionsWithLockId("test-mongo-lock-client-opt")
	lock2, err := NewMongoLockByClientWithOptions(ctx, client, options)
	assert.Nil(t, err)
	assert.NotNil(t, lock2)

	assert.Nil(t, lock.Lock(ctx, "owner-1"))
	assert.Nil(t, lock.UnLock(ctx, "owner-1"))
}

func TestMongoLockByClient_Errors(t *testing.T) {
	old := newMongoStorageByClient
	newMongoStorageByClient = func(ctx context.Context, options *mongodb_storage.MongoStorageOptions) (*mongodb_storage.MongoStorage, error) {
		return nil, errors.New("boom")
	}
	defer func() { newMongoStorageByClient = old }()

	testMongoURI(t)
	ctx := context.Background()
	// 用独立 client 避免缓存污染
	ctx2, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	boomClient, err := mongo.Connect(ctx2, options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	assert.Nil(t, err)
	defer boomClient.Disconnect(context.Background())

	_, err = GetMongoLockFactoryByClient(ctx, boomClient)
	assert.NotNil(t, err)
	_, err = NewMongoLockByClient(ctx, boomClient, "x")
	assert.NotNil(t, err)
	_, err = NewMongoLockByClientWithOptions(ctx, boomClient, storage_lock.NewStorageLockOptionsWithLockId("x"))
	assert.NotNil(t, err)
}

func TestNewMongoLockByUri_Local(t *testing.T) {
	uri := testMongoURI(t)
	ctx := context.Background()

	factory, err := GetMongoLockFactoryByUri(ctx, uri)
	assert.Nil(t, err)
	assert.NotNil(t, factory)

	// 缓存命中
	factory2, err := GetMongoLockFactoryByUri(ctx, uri)
	assert.Nil(t, err)
	assert.Equal(t, factory, factory2)

	lock, err := NewMongoLockByUri(ctx, uri, "test-mongo-lock-uri")
	assert.Nil(t, err)
	assert.NotNil(t, lock)

	opts := storage_lock.NewStorageLockOptionsWithLockId("test-mongo-lock-uri-opt")
	lock2, err := NewMongoLockByUriWithOptions(ctx, uri, opts)
	assert.Nil(t, err)
	assert.NotNil(t, lock2)

	assert.Nil(t, lock.Lock(ctx, "owner-1"))
	assert.Nil(t, lock.UnLock(ctx, "owner-1"))
}

func TestMongoLockByUri_InjectError(t *testing.T) {
	old := newMongoStorageByURI
	newMongoStorageByURI = func(ctx context.Context, options *mongodb_storage.MongoStorageOptions) (*mongodb_storage.MongoStorage, error) {
		return nil, errors.New("boom")
	}
	defer func() { newMongoStorageByURI = old }()

	ctx := context.Background()
	uri := "mongodb://127.0.0.1:27017-inject-error"
	_, err := GetMongoLockFactoryByUri(ctx, uri)
	assert.NotNil(t, err)
	_, err = NewMongoLockByUri(ctx, uri, "x")
	assert.NotNil(t, err)
	_, err = NewMongoLockByUriWithOptions(ctx, uri, storage_lock.NewStorageLockOptionsWithLockId("x"))
	assert.NotNil(t, err)
}
