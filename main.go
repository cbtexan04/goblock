package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func displayChain(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func writeBlock(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func run() error {
	mr := mux.NewRouter()
	mr.HandleFunc("/", displayChain).Methods("GET")
	mr.HandleFunc("/", writeBlock).Methods("POST")

	port := os.Getenv("ADDR")

	s := &http.Server{
		Addr:           ":" + port,
		Handler:        mr,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Listening on ", port)
	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
