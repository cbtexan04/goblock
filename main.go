package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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
	Write(w, 200, Blockchain)
	bcMutex.Unlock()
}

type Message struct {
	BPM int `json:"bpm"`
}

func writeBlock(w http.ResponseWriter, r *http.Request) {
	var m Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	// This is OK, since we've established our genesis block already
	lastBlock := Blockchain[len(Blockchain)-1]

	newBlock, err := lastBlock.Generate(m.BPM)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if lastBlock.IsValidNextBlock(newBlock) {
		newBlockchain := append(Blockchain, newBlock)

		// We could run into an issue where two nodes have added
		// (valid) blocks to their chains and we get them both. In this
		// case, we need to make sure to pick the longest chain.
		if len(newBlockchain) > len(Blockchain) {
			Blockchain = newBlockchain
		}
	}

	Write(w, http.StatusCreated, newBlock)
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

	log.Printf("Listening on %s", port)
	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func main() {
	// Load our port through the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Create our genesis block
	t := time.Now()
	genesis := &Block{0, t.String(), "", "", 0}
	Blockchain = append(Blockchain, genesis)

	log.Fatal(run())
}
