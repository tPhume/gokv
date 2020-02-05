package main

import (
	"context"
	"github.com/tPhume/gokv/kv"
	"google.golang.org/grpc"
	"log"
)

// Build this and start the grpcServer
// connects to the default grpcServer
func main() {
	testValue := map[string]string{"greeting": "Hello, I am test value!"}

	// attempt to connect to the grpcServer
	conn, err := grpc.Dial("127.0.0.1:9999", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to grpcServer, %s", err)
	}
	defer conn.Close()

	// pass the connection (channel) to create the client
	client := kv.NewGoKvClient(conn)

	// insert
	response, err := client.Insert(context.Background(), &kv.KeyValue{
		Key:   &kv.Key{Key: "Test"},
		Value: &kv.Value{Value: testValue},
	})

	if err != nil {
		log.Fatalf("insert failed, %s", err)
	} else {
		log.Println(response)
	}
	
	// search
	response, err = client.Search(context.Background(), &kv.Key{Key: "Test"})
	if err != nil {
		log.Fatalf("search failed, %s", err)
	} else {
		log.Println(response)
	}
}
