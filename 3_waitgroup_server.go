package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	// we don’t see our logs saying each of our servers are closed.
	// We need a way for the goroutines to signal to the main function they are done and for the main function to wait until they are all done.

	// Using waitgroup, we can say to the main to wait until all go routine are done
	var wg sync.WaitGroup
	wg.Add(1) // we have 2 go routines

	ctx, cancel := context.WithCancel(context.Background())

	go newHelloWorldWaitGracefulServer(ctx, &wg)
	go newHowAreYouWaitGracefulServer(ctx, &wg)

	fmt.Println("All servers are started")

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-signals
	fmt.Println("Terminating ")
	cancel()

	wg.Wait() // we need to wait that all go routines are done to finish
}

func newHelloWorldWaitGracefulServer(ctx context.Context, wg *sync.WaitGroup) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`Hello, world!`))
	})

	server := &http.Server{Addr: ":7000", Handler: mux}

	go func() {
		<-ctx.Done()
		shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(shutCtx); err != nil {
			fmt.Printf("error shutting down the hello world server: %s\n", err)
		}
		fmt.Println("the hello world server is closed")
		wg.Done() // this go routine is done and removed from the group
	}()

	// this part must go after the go routine as it's blocking
	fmt.Println("the hello world server is starting")
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Printf("error starting the hello world server: %s\n", err)
	}

	fmt.Println("the hello world server is closing")
}

func newHowAreYouWaitGracefulServer(ctx context.Context, wg *sync.WaitGroup) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("How are you?"))
	})

	server := &http.Server{Addr: ":8000", Handler: mux}

	go func() {
		<-ctx.Done()
		shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(shutCtx); err != nil {
			fmt.Printf("error shutting down the how are you server: %s\n", err)
		}
		fmt.Println("the how are you server is closed")
		wg.Done() // this go routine is done and removed from the group
	}()

	fmt.Println("the how are you server is starting")
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Printf("error starting the how are you server: %s\n", err)
	}

	fmt.Println("the how are you server is closing")
}