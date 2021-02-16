package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	helloWorldSvr := newHelloWorldServer()
	helloNameSvr := newHowAreYouServer()

	// 1- the ListenAndServe method on an http.Server struct is a blocking call, so the program stops after this line
	// 2- trying to add go routine
	go helloWorldSvr.ListenAndServe()
	go helloNameSvr.ListenAndServe()

	fmt.Println("All servers are started")

	// 2 - with go routines, main finishes first. How can we make the program wait to complete until we give it the signal to quit?

	// 3 - signals is a channel of size 1
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals // 3- This line is asking for a value off of the channel. This is a blocking operation that will only continue when it gets a value or the channel is closed.

	// 3 -Therefore, this line prevents our program from continuing until it receives a SIGINT or SIGTERM, which would be sent by us when we stop the program

	// 4- Make the servers gracefully shutdown

}

func newHelloWorldServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`Hello, world!`))
	})

	return &http.Server{Addr: ":7000", Handler: mux}
}

func newHowAreYouServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("How are you?"))
	})

	return &http.Server{Addr: ":8000", Handler: mux}
}
