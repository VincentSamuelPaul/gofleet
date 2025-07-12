package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

var mtx sync.Mutex
var wg sync.WaitGroup

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

func checkServerStatus(w http.ResponseWriter, r *http.Request) error {
	status := Status()
	return WriteJSON(w, http.StatusOK, map[string]any{"nodes": status})
}

var nodesStatus map[string]bool

func work(w http.ResponseWriter, r *http.Request) error {
	return router(w, r)
}

func router(w http.ResponseWriter, r *http.Request) error {
	freeNode := GetFreeNode(nodeStatus)
	if freeNode == "" {
		return WriteJSON(w, http.StatusServiceUnavailable, ErrorMsg{Error: "No available backend nodes"})
	}

	mtx.Lock()
	nodeStatus[freeNode] = false
	mtx.Unlock()

	res, err := http.Get(freeNode + "/work")
	if err != nil {
		mtx.Lock()
		nodeStatus[freeNode] = true
		mtx.Unlock()

		return WriteJSON(w, http.StatusBadGateway, ErrorMsg{Error: err.Error()})
	}
	defer res.Body.Close()

	var req WorkOutput
	if err = json.NewDecoder(res.Body).Decode(&req); err != nil {
		return WriteJSON(w, http.StatusInternalServerError, ErrorMsg{Error: err.Error()})
	}

	mtx.Lock()
	nodeStatus[freeNode] = true
	mtx.Unlock()
	fmt.Println("Recieved from node: " + freeNode + " | Prime chosen: " + fmt.Sprintf("%d", req.Number))
	return WriteJSON(w, http.StatusOK, req)
}

func work1(w http.ResponseWriter, r *http.Request) error {
	node := GetBackendNode()
	res, err := http.Get(node + "/work")
	if err != nil {
		return WriteJSON(w, http.StatusBadGateway, ErrorMsg{Error: err.Error()})
	}
	defer res.Body.Close()

	var req WorkOutput
	if err = json.NewDecoder(res.Body).Decode(&req); err != nil {
		return WriteJSON(w, http.StatusInternalServerError, ErrorMsg{Error: err.Error()})
	}
	fmt.Println("Recieved from node: " + node + " | Prime chosen: " + fmt.Sprintf("%d", req.Number))
	return WriteJSON(w, http.StatusOK, req)
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/status", makeHTTPHandleFunc(checkServerStatus))
	router.HandleFunc("/work", makeHTTPHandleFunc(work1))
	log.Printf("\n\nLoad Balancer running on port: %+v\n\n", s.listenAddr)
	nodesStatus = Status()
	for key, value := range nodesStatus {
		fmt.Println("node:", key, "| running:", value)
	}
	fmt.Println("\nRouter Logs")
	http.ListenAndServe(s.listenAddr, router)
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
