package main

import (
	"context"
	"log"
	"logistics-simulator/internal/database"
	"logistics-simulator/internal/handlers"
	"logistics-simulator/internal/workers"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Println("--- Khởi động Wolt Simulator ---")

	// 1. init database connection and auto-migrate Order model
	database.InitDB()

	// 2. Create a channel to send order IDs from API handlers to Workers
	jobChan := make(chan uint, 100)

	// 3. run worker pool with 3 workers to process orders concurrently
	workers.StartWorkerPool(3, jobChan, database.DB)

	// 4. init OrderHandler with DB and JobChan
	orderHandler := &handlers.OrderHandler{
		DB:      database.DB,
		JobChan: jobChan,
	}

	// 5. Setup Router with OrderHandler
	router := handlers.SetupRouter(orderHandler)

	// 6. start HTTP server in a separate Goroutine
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// 7. run server in Goroutine so that it doesn't block the main thread, allowing us to listen for shutdown signals
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server not working: %s\n", err)
		}
	}()

	// 8.Graceful Shutdown: Listen for OS signals to gracefully shut down the server and workers
	// create channel to listen for interrupt or terminate signals
	quit := make(chan os.Signal, 1)
	// listen for SIGINT (Ctrl+C) and SIGTERM (termination signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// await signal in main thread, when received, proceed to shutdown
	<-quit
	log.Println("shutdown signal received, shutting down server...")

	// 9. set timeout context for server shutdown to ensure it doesn't hang indefinitely
	// set timeout 15 seconds to allow ongoing requests to complete and workers to finish processing current orders
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close server, if error occurs during shutdown, log it
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown error:", err)
	}

	// close job channel to signal workers to stop after finishing current jobs
	close(jobChan)

	log.Println("server gracefully stopped, goodbye! 👋")
}
