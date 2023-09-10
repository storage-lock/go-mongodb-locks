package main

import (
	"context"
	"fmt"
	mongodb_locks "github.com/storage-lock/go-mongodb-locks"
	storage_lock "github.com/storage-lock/go-storage-lock"
)

func main() {

	// 第一步：从mongo uri创建一把锁
	mongoUri := "mongodb://root:UeGqAm8CxYGldMDLoNNt@127.0.0.1:27017/?connectTimeoutMS=300000"
	lock, err := mongodb_locks.NewMongoLockByUri(context.Background(), mongoUri, "foo-lock")
	if err != nil {
		panic(err)
	}
	// 第二步：获取锁
	ownerA := "owner-A"
	err = lock.Lock(context.Background(), ownerA)
	if err != nil {
		panic(err)
	}
	// 第三步：一定要记得释放锁，否则分布式系统中的其他人将无法获取锁
	defer func(lock *storage_lock.StorageLock, ctx context.Context, ownerId string) {
		err := lock.UnLock(ctx, ownerId)
		if err != nil {
			panic(err)
		}
	}(lock, context.Background(), ownerA)

	// 第四步：接下来直到函数结束都属于临界区，可以在这些地方对资源进行操作
	fmt.Println("资源操作")

}
