package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Abir-Zayn/studentsApi/internal/config"
	"github.com/Abir-Zayn/studentsApi/internal/http/handlers/student"
	"github.com/Abir-Zayn/studentsApi/internal/storage/sqlite"
)

func main() {
	// load config
	cfg := config.MustLoad()
	// database setup

	store, err := sqlite.New(*cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer store.Db.Close()


	// setup router
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New(store))

	// setup http server
	server := http.Server {
		Addr : cfg.HTTPServer.Addr, //where to listen?
		Handler: router,  // which router to use?
	}
	slog.Info("Starting server", slog.String("address: ", cfg.HTTPServer.Addr))
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Signal >> A message sent to a process by the operating system or another process to notify it of an event.
	// SIGINT >> A signal sent to a process by the operating system when the user wishes to interrupt the process.
	// SIGTERM >> A signal sent to a process by the operating system to request its termination.

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}()
	<- done

	slog.Info("Shutting Down server")

	// GoRoutine >> [SUPER lightweight thread managed by the Go runtime]
	//A goroutine means multiple tasks can be in progress at the same time.

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}else {
		slog.Info("Server exited properly")
	}

	
} 