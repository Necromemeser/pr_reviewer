package main

import (
	"log"
	"net/http"
	"pr_reviewer/internal/handlers"
	"pr_reviewer/internal/storage"
)

func main() {
	db, err := storage.ConnectAndMigrate()
	if err != nil {
		log.Fatalf("DB error: %v", err)
	}
	defer db.Close()

	st := storage.NewStorage(db)
	mux := handlers.NewRouter(st)

	log.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", mux)
}
