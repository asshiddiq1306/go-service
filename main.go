package main

import (
	"context"
	"go-service/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	l := log.New(os.Stdout, "product-api ", log.LstdFlags)
	ph := handlers.NewProduct(l)

	sm := http.NewServeMux()
	sm.Handle("/", ph)

	s := http.Server{
		Addr:         ":8080",
		Handler:      sm,
		ErrorLog:     l,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sigChan := <-c
	l.Printf("Receice server shutdown gracefully %s\n", sigChan)

	tc, _ := context.WithTimeout(context.Background(), 25*time.Second)
	s.Shutdown(tc)
}
