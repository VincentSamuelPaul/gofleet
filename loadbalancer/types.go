package main

import "time"

type APIServer struct {
	listenAddr string
}

type ErrorMsg struct {
	Error string `json:"error"`
}

type NodeStatus struct {
	Alive bool `json:"alive"`
}

type WorkOutput struct {
	Number  int           `json:"random_number"`
	Factors []int         `json:"prime_factors"`
	Time    time.Duration `json:"execution_time"`
}
