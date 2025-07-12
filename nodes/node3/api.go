package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

var workQueue = make(chan WorkRequest, 100)

func startWorkers() {
	for i := 0; i < 4; i++ {
		go worker(i)
	}
}

func worker(id int) {
	for job := range workQueue {
		fmt.Printf("Worker %d: Started job\n", id)

		result := FindPrimes()
		job.Result <- result
		fmt.Printf("Worker %d: Finished job\n", id)
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	startWorkers()
	router.HandleFunc("/status", makeHTTPHandleFunc(checkServerStatus))
	router.HandleFunc("/work", makeHTTPHandleFunc(work))
	log.Println("\nJSON API server running on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func checkServerStatus(w http.ResponseWriter, r *http.Request) error {
	return WriteJSON(w, http.StatusOK, map[string]bool{"alive": true})
}

func work(w http.ResponseWriter, r *http.Request) error {

	resultChan := make(chan WorkOutput)

	job := WorkRequest{
		Result: resultChan,
	}

	select {
	case workQueue <- job:
		result := <-resultChan
		WriteJSON(w, http.StatusOK, result)

	default:
		WriteJSON(w, http.StatusBadGateway, ErrorMsg{Error: "Server busy"})
	}

	return nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ErrorMsg{Error: err.Error()})
		}
	}
}
