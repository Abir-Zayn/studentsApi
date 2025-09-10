package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Abir-Zayn/studentsApi/internal/config"
)

func main() {
	// load config

	cfg := config.MustLoad()

	// database setup
	// setup router
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, World!"))
		})

	// setup server
	server := http.Server {
		Addr : cfg.Addr,
		Handler: router,
	}

	fmt.Printf("Starting server on %s", cfg.HTTPServer.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
	fmt.Println("Server started on", cfg.Addr)
} 