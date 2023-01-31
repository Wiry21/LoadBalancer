package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"sync"
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

		resp, err := client.R().
			EnableTrace().
			Get("http://localhost:8080")

		//Explore response object
		if resp.StatusCode() != 200 {
			fmt.Println("Response Info:")
			fmt.Println("  Error      :", err)
			fmt.Println("  Status Code:", resp.StatusCode())
		}
	}
}
