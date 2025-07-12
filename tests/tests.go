package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

var lb string = "http://localhost:3000/work"

var wg sync.WaitGroup

func main() {
	mainStart := time.Now()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			start := time.Now()
			res, err := http.Get(lb)
			duration := time.Since(start)

			if err != nil {
				fmt.Println(err)
			}
			defer res.Body.Close()

			fmt.Printf("Request id: %d | ExecTime: %+v Prime: %+v\n", i, duration, res.Body)
		}(i)
	}
	wg.Wait()
	finalDuration := time.Since(mainStart)
	fmt.Printf("\nTotal execution time: %+v\n", finalDuration)
}
