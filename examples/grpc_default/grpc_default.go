package main

import (
	"github.com/tPhume/gokv/kv"
	"log"
	"net"
)

// Build this and run it to open the grpc server
// clients can now connect to it via localhost:8080
func main() {
	addr := "0.0.0.0:9999"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("cannot listen on port, %s", err)
	}

	grpcServer := kv.DefaultGrpcServer()
	log.Fatal(grpcServer.Serve(lis))
}
