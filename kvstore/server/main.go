package main

import (
	"context"
	"net"
	"sync"
	"time"

	kvcache "github.com/VincentSamuelPaul/gofleet/kvstore/cache"
	kv "github.com/VincentSamuelPaul/gofleet/kvstore/kvstore/proto"

	"log"

	"google.golang.org/grpc"
)

var mtx sync.RWMutex

type server struct {
	kv.UnimplementedKVStoreServer
}

var cacheCapacity = 5

var store = map[string]string{}
var cache = kvcache.NewLRUCache(cacheCapacity)

func (s *server) Get(ctx context.Context, req *kv.Key) (*kv.Value, error) {
	now := time.Now()
	value, found := cache.Get(req.Key)
	if found {
		log.Printf("Get() key=%s found=%t\n execution-time=%s", req.Key, found, time.Since(now))
		return &kv.Value{Value: value, Found: true}, nil
	}

	mtx.RLock()
	value, found = store[req.Key]
	mtx.RUnlock()

	if !found {
		log.Println("Get() no key found")
		return &kv.Value{Value: "", Found: false}, nil
	}
	cache.Set(req.Key, value)
	log.Printf("Get() loaded key=%s into cache from kvstore, execution-time=%s", req.Key, time.Since(now))
	return &kv.Value{Value: value, Found: true}, nil
}

func (s *server) Set(ctx context.Context, req *kv.KeyValue) (*kv.Ack, error) {
	mtx.Lock()
	store[req.Key] = req.Value
	mtx.Unlock()

	cache.Set(req.Key, req.Value)
	log.Printf("Set() key=%s value=%s", req.Key, req.Value)
	return &kv.Ack{Ok: true}, nil
}

func main() {
	const port = ":3000"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("gRPC server listening on port %s\n", port)
	grpcServer := grpc.NewServer()
	kv.RegisterKVStoreServer(grpcServer, &server{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve GRPC server: %v", err)
	}
}
