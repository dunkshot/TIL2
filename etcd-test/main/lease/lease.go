package main

import (
	"context"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {

	// Watcher Section
	// go func() {
	// 	watcherClient, err := clientv3.New(clientv3.Config{
	// 		Endpoints:   []string{"localhost:22379"},
	// 		DialTimeout: 5 * time.Second,
	// 	})
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	defer watcherClient.Close()

	// 	rch := watcherClient.Watch(context.Background(), "fooForLease")
	// 	for wresp := range rch {
	// 		for _, ev := range wresp.Events {
	// 			fmt.Printf("Watch - %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
	// 		}
	// 	}
	// }()

	// Lease Section
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:22379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// minimum lease TTL is 5-second
	resp, err := cli.Grant(context.TODO(), 5)
	if err != nil {
		log.Fatal(err)
	}

	// after 5 seconds, the key 'fooForLease' will be removed
	_, err = cli.Put(context.TODO(), "fooForLease", "bar", clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatal(err)
	}
	println(resp.ID)

	// renew the lease only once
	ka, kaerr := cli.KeepAliveOnce(context.TODO(), resp.ID)
	if kaerr != nil {
		log.Fatal(kaerr)
	}
	fmt.Println("ttl: ", ka.TTL)
}
