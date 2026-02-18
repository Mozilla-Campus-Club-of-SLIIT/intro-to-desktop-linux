package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Mozilla-Campus-Club-of-SLIIT/intro-to-desktop-linux/internal/engine"
)

func main() {
	dbAddr := os.Getenv("DB_ADDR")
	if dbAddr == "" {
		dbAddr = "localhost:6379" // Default for local development
	}

	db, err := engine.NewDB(dbAddr)
	if err != nil {
		log.Fatalf("Failed to connect to Valkey: %v", err)
	}
	defer db.Close()

	log.Printf("Connected to Valkey at %s", dbAddr)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Valkey API Server is running"))
	})

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
