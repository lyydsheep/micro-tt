package biz

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"testing"
	"time"
)

func TestEtcd(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:22379"},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()

	ctx := context.TODO()
	if _, err = cli.Put(ctx, "key", "value"); err != nil {
		t.Error(err)
	}
	t.Log("put ok")

	resp, err := cli.Get(ctx, "key")
	if err != nil {
		t.Error(err)
	}
	for _, kv := range resp.Kvs {
		t.Log(string(kv.Key), string(kv.Value))
	}

	if _, err = cli.Delete(ctx, "key"); err != nil {
		t.Error(err)
	}
}

func TestLease(t *testing.T) {
	// 续上文的 main 函数
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:22379"},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		t.Fatal(err)
	}
	// 1. 创建租约（10秒过期）
	lease, err := cli.Grant(context.TODO(), 1)
	if err != nil {
		t.Fatal(err)
	}

	// 2. 绑定键值到租约
	_, err = cli.Put(context.TODO(), "/demo/lease_key", "lease_value", clientv3.WithLease(lease.ID))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("键值绑定租约成功")

	// 3. 自动续约（每5秒续一次）
	keepAlive, err := cli.KeepAlive(context.TODO(), lease.ID)
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for resp := range keepAlive {
			fmt.Printf("租约续约成功: ID=%x, TTL=%d, now: %v\n", resp.ID, resp.TTL, time.Now().Format("2006-01-02 15:04:05"))
		}
	}()

	// 4. 等待租约过期
	time.Sleep(15 * time.Second)
	resp, err := cli.Get(context.TODO(), "/demo/lease_key")
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Kvs) == 0 {
		fmt.Println("键值已随租约过期自动删除")
	}
	for _, kv := range resp.Kvs {
		fmt.Println(string(kv.Key), string(kv.Value))
	}
}
