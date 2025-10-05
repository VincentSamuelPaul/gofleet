package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/VincentSamuelPaul/gofleet/kvstore/kvstore/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewKVStoreClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("\nPopulating store with keys...")
	for i := 1; i <= 10; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		_, err := client.Set(ctx, &pb.KeyValue{Key: key, Value: value})
		if err != nil {
			log.Fatalf("Set error: %v", err)
		}
	}

	fmt.Println("Populating cache...")
	for i := 1; i <= 5; i++ {
		key := fmt.Sprintf("key%d", i)
		_, err := client.Get(ctx, &pb.Key{Key: key})
		if err != nil {
			log.Fatalf("Get error: %v", err)
		}
	}

	fmt.Println("\nTesting execution times (cache hit vs miss):")
	testKeys := []string{"key3", "key7"}
	for _, key := range testKeys {
		start := time.Now()
		resp, err := client.Get(ctx, &pb.Key{Key: key})
		if err != nil {
			log.Fatalf("Get error: %v", err)
		}
		fmt.Printf("\nKey=%s, Found=%t, Value=%s, ExecutionTime=%s\n\n",
			key, resp.Found, resp.Value, time.Since(start))
	}
}
