package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func (s server) setupParams() error {
	dt, err := GetLayout("graph.json")
	// log.Printf("setup: Graph Parameters %v\n", dt)
	return dt, err
}

// func (s server) setupData() ( error) {
// 	var cldt []ClientData
// 	cldt, err := server.GetClientData("/data/datapoints.json")
// 	log.Printf("setup: Datapoints %v\n", cldt)

// 	return &cldt, err
// }

// Main wraps the run to allow a start failure
func main() {
	var addr = flag.String("addr", ":8080", "The address of the server.")
	flag.Parse()

	log.Println("Starting webserver on", *addr)

	var s Server
	err := s.setupParams()
	if err != nil {
		fmt.Printf("%#v\n", err)
	}

	// err := s.setupData()
	// if err != nil {
	// 	fmt.Printf("%#v\n", err)
	// }

	log.Printf("main: dt %v\n", dt)
	log.Printf("main: cldt %v\n", cldt)

	http.HandleFunc("/", dt.TestHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
