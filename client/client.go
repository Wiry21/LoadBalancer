package main

import (
	"sync"

	"github.com/go-resty/resty/v2"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(25)

	for i := 0; i < 25; i++ {
		go func() {
			defer wg.Done()
			heartBeat()
		}()
	}
	wg.Wait()
}

func heartBeat() {
	for i := 0; i < 250; i++ {
		//Create a Resty Client
		client := resty.New()

		_, _ = client.R().
			EnableTrace().
			Get("http://host.docker.internal:8080")
	}
}
