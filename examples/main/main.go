package main

import (
	"github.com/tPhume/gokv/btree"
	"github.com/tPhume/gokv/kv"
	"log"
	"net"
)

func main() {
	restAddr := "0.0.0.0:8888"
	grpcAddr := "0.0.0.0:9999"

	store := btree.NewBtree(3)

	restServer := kv.RestWithStore(store)
	grpcServer := kv.GrpcWithStore(store)

	go func() {
		log.Fatal(restServer.Run(restAddr))
	}()

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("could not listen on port %s", err)
	}

	log.Fatal(grpcServer.Serve(lis))
}
