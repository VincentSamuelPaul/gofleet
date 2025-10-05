package main

import (
	"context"
	"log"
	"net"
	"testing"

	kvstore "github.com/VincentSamuelPaul/gofleet/kvstore/kvstore/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	kvstore.UnimplementedKVStoreServer
}

// var server kvstore.UnimplementedKVStoreServer

func startTestServer(t *testing.T) (*grpc.Server, kvstore.KVStoreClient, func()) {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	s := grpc.NewServer()
	kvstore.RegisterKVStoreServer(s, &server{})

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Printf("server exited: %v", err)
		}
	}()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}

	client := kvstore.NewKVStoreClient(conn)

	cleanup := func() {
		s.Stop()
		conn.Close()
		lis.Close()
	}
	return s, client, cleanup
}

func TestKVStore_GetSet(t *testing.T) {
	_, client, cleanup := startTestServer(t)
	defer cleanup()

	ctx := context.Background()

	_, err := client.Set(ctx, &kvstore.KeyValue{Key: "key1", Value: "value1"})
	if err != nil {
		t.Fatalf("Set() failed: %v", err)
	}

	resp, err := client.Get(ctx, &kvstore.Key{Key: "key1"})
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}
	if !resp.Found || resp.Value != "value1" {
		t.Errorf("Expected value='value1', got %v (found=%v)", resp.Value, resp.Found)
	}
}
