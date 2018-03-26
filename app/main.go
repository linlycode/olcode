package main

import (
	"log"

	"olcode"
)

func main() {
	var port int16 = 5432
	var homePath = "/Users/xkwei/github-workspace/yungewu-cloudapp_go/src/olcode/client/dist"
	s := olcode.NewService(port, homePath)

	log.Printf("starting service at port %d", port)
	if err := s.Serve(); err != nil {
		log.Printf("failed to start service, %s", err)
	}
}
