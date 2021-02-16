package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// create context
	ctx, cancel := context.WithCancel(context.Background())

	go newHelloWorldGracefulServer(ctx)
	go newHowAreYouGracefulServer(ctx)

	// 1- We can see that we tried to start a server that was already closed.
	// This happened because the goroutine in each server function runs before the server is even started.

	// 2- Contexts are used all throughout the Golang standard library as a way to signal when to stop things
	// from running and also as a way to pass values down to functions

	fmt.Println("All servers are started")

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	// We want to cancel the context once we get a signal to stop everything.
	//So let’s cancel the context right after we get that signal!
	cancel()
}

func newHelloWorldGracefulServer(ctx context.Context) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`Hello, world!`))
	})

	server := &http.Server{Addr: ":7000", Handler: mux}

	// go routine to close gracefully the server
	// 1- Runs before the server could start, we need to find a way to run that only on shutdown
	go func() {
		fmt.Print("Waiting to shutdown")
		// 2- This can be done by waiting for a message on the Done channel of the context.
		// This channel will be closed when we call cancel from our main function.
		// This means any places we are waiting for a message will continue moving forward.
		<-ctx.Done()
		shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		fmt.Print("Shutting down")
		if err := server.Shutdown(shutCtx); err != nil {
			fmt.Printf("error shutting down the hello world server: %s\n", err)
		}
		fmt.Println("the hello world server is closed")
	}()

	// starting server, the ListenAndServe is blocking and must go after the go routine
	fmt.Println("the hello world server is starting")
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Printf("error starting the hello world server: %s\n", err)
	}

	fmt.Println("the hello world server is closing")
}

func newHowAreYouGracefulServer(ctx context.Context) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("How are you?"))
	})

	server := &http.Server{Addr: ":8000", Handler: mux}

	// go routine to close gracefully the server
	// 1- Runs before the server could start, we need to find a way to run that only on shutdown
	go func() {
		// 2- This can be done by waiting for a message on the Done channel of the context.
		// This channel will be closed when we call cancel from our main function.
		// This means any places we are waiting for a message will continue moving forward.
		<-ctx.Done()
		shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(shutCtx); err != nil {
			fmt.Printf("error shutting down the how are you server: %s\n", err)
		}
		fmt.Println("the how are you server is closed")
	}()

	// starting server, the ListenAndServe is blocking and must go after the go routine
	fmt.Println("the how are you server is starting")
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Printf("error starting the how are you server: %s\n", err)
	}

	fmt.Println("the how are you server is closing")
}