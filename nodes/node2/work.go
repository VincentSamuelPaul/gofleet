package main

import (
	"math/rand"
	"sync"
	"time"
)

func FindPrimes() WorkOutput {
	n := rand.Intn(100000)
	var mtx sync.Mutex
	var wg sync.WaitGroup
	var factors []int
	t := time.Now()
	for i := 2; i <= n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var status bool = true
			for j := 2; j <= i; j++ {
				if i == j {
					continue
				} else if i%j == 0 {
					status = false
					break
				}
			}
			if status {
				mtx.Lock()
				factors = append(factors, i)
				mtx.Unlock()
			}
		}(i)
	}
	wg.Wait()
	return WorkOutput{Number: n, Factors: factors, Time: time.Duration(time.Since(t).Microseconds())}
}
