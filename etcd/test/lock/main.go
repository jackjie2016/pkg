package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"sync"
	"time"
)

var sg sync.WaitGroup

type Counter struct {
	count int
}

func (m *Counter) Incr() {
	m.count++
}

func (m *Counter) Count() int {
	return m.count
}

func main() {
	endpoints := []string{"http://10.10.4.42:12379", "http://10.10.4.42:22379", "http://10.10.4.42:32379"}
	client, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	counter := &Counter{}

	session, err := concurrency.NewSession(client)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// 此处使用newMutex初始化
	locker := concurrency.NewMutex(session, "/my-test-lock3")
	fmt.Println("locking...", time.Now().Format("2006-01-02 15:04:05"))
	err = locker.TryLock(context.Background())
	// 获取锁失败就抛错
	if err != nil {
		fmt.Println("lock failed", err)
		return
	}
	fmt.Println("locked...", time.Now().Format("2006-01-02 15:04:05"))
	time.Sleep(5 * time.Second)
	counter.Incr()
	err = locker.Unlock(context.Background())
	if err != nil {
		fmt.Println("unlock failed", err)
		return
	}
	fmt.Println("released...", time.Now().Format("2006-01-02 15:04:05"))

	fmt.Println("count:", counter.Count())
}
