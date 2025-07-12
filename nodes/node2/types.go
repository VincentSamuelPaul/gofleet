package main

import "time"

type ErrorMsg struct {
	Error string `json:"error"`
}

type WorkOutput struct {
	Number  int           `json:"random_number"`
	Factors []int         `json:"prime_factors"`
	Time    time.Duration `json:"execution_time"`
}

type WorkRequest struct {
	Result chan WorkOutput
}
