package main

import (
	"github.com/tPhume/gokv/kv"
	"log"
)

func main() {
	router := kv.DefaultRestServer()
	log.Fatal(router.Run("0.0.0.0:8888"))
}
