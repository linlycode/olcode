package main

import (
	"flag"
	"log"

	"olcode"
)

var port = flag.String("port", "5432", "service port")
var staticPath = flag.String("static_path", "$GOPATH/src/olcode/client/dist", "static files dir")

func main() {
	flag.Parse()
	s := olcode.NewService(*port, *staticPath)

	log.Printf("starting service at port %s", *port)
	if err := s.Serve(); err != nil {
		log.Printf("failed to start service, %s", err)
	}
}
