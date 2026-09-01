package main

import (
	"log"
	"os"
)

func main() {
	cfg := config{
		addr: ":8080",
		db:   dbConfig{},
	}

	api := application{
		config: cfg,
	}

	handler := api.mount()
	if err := api.run(&handler); err != nil {
		log.Printf("Error starting server: %v", err)
		os.Exit(1)
	}
}
