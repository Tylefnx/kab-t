package main

import (
	"game-server/server"
	"log"
	"net/http"
)

func main() {
	srv := server.NewServer(2) // Örnek max kullanıcı sayısı

	corsHandler := func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == "OPTIONS" {
				return
			}
			handler.ServeHTTP(w, r)
		})
	}

	http.HandleFunc("/ws", srv.Handler().ServeWs)
	http.HandleFunc("/answer", srv.Handler().HandleAnswer)
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler(http.DefaultServeMux)))
}
