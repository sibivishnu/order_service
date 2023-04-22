package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/sibivishnu/order_service/handlers"
)

// Run starts the HTTP server and sets up graceful shutdown.
func Run() {
	// Create a new router
	router := mux.NewRouter()

	// Register the order-related handlers
	router.HandleFunc("/api/orders", handlers.AddOrder).Methods("POST")
	router.HandleFunc("/api/orders/{id}", handlers.UpdateOrder).Methods("PUT")
	router.HandleFunc("/api/orders", handlers.GetOrders).Methods("GET")

	// Set up the HTTP server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start the server
	go func() {
		fmt.Println("Starting server on port 8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Failed to start server: %s\n", err)
			os.Exit(1)
		}
	}()

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// Shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Failed to gracefully shut down server: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Server shut down")
}
