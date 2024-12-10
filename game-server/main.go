package main

import (
	"game-server/server"
	"log"
)

func main() {
	srv := server.NewServer(5) // Örnek max kullanıcı sayısı
	log.Fatal(srv.ListenAndServe(":8080"))
}
