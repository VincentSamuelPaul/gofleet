package main

import (
	"encoding/json"
	"net/http"
)

var nodeStatus = map[string]bool{"http://localhost:3001": false, "http://localhost:3002": false, "http://localhost:3003": false}

func Status() map[string]bool {
	for key, _ := range nodeStatus {
		res, err := http.Get(key + "/status")
		if err != nil {
			continue
		}
		var req NodeStatus
		if err = json.NewDecoder(res.Body).Decode(&req); err != nil {
			continue
		}
		if req.Alive {
			nodeStatus[key] = true
		}
	}
	return nodeStatus
}

func GetFreeNode(nodes map[string]bool) string {
	val := ""
	for key, value := range nodes {
		if value == true {
			val = key
			break
		}
	}
	return val
}

var backends = []string{
	"http://localhost:3001",
	"http://localhost:3002",
	"http://localhost:3003",
}

var backendIndex int

func GetBackendNode() string {
	mtx.Lock()
	defer mtx.Unlock()
	node := backends[backendIndex]
	backendIndex = (backendIndex + 1) % len(backends)
	return node
}
