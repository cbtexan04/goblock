package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func Write(w http.ResponseWriter, code int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		log.Println(err)
		WriteError(w, 500, "unable to marshal json response")
		return
	}

	w.WriteHeader(code)
	w.Write(b)
}

func WriteError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	data := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{
		"failed", message,
	}

	b, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "something has gone horribly wrong", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	w.Write(b)
}

func displayChain(w http.ResponseWriter, r *http.Request) {
	// While locking the blockchain isn't stricly necessary at this point,
	// it's probably a good idea
	bcMutex.Lock()
	Write(Blockchain)
	bcMutex.Unlock()
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
