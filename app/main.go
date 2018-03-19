package main

import (
	"log"

	"olcode"
)

func main() {
	var port int16 = 5432
	s := olcode.NewService(port)

	log.Printf("starting service at port %d", port)
	if err := s.Serve(); err != nil {
		log.Printf("failed to start service, %s", err)
	}
}
