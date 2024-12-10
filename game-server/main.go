package main

import (
	"game-server/server"
	"log"
)

func main() {
	srv := server.NewServer(5)
	log.Fatal(srv.ListenAndServe(":8080"))
}
