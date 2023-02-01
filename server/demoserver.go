package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

var port = flag.Int("port", 8082, "Port to start the demo service on")

type DemoServer struct{}

func (f *DemoServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Second * 1)
	res.WriteHeader(200)
	res.Write([]byte(fmt.Sprintf("All good! from server %d.", *port)))
}

func main() {
	flag.Parse()

	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), &DemoServer{}); err != nil {
		log.Fatal(err)
	}
}
