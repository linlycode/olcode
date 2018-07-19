package main

import (
	"flag"
	"log"

	"github.com/linlycode/olcode/pkg/api"
)

var port = flag.String("port", "8081", "service port")

func main() {
	flag.Parse()
	log.Printf("serve on port %s", *port)
	s := api.NewService(*port)
	if err := s.Serve(); err != nil {
		log.Printf("server stopped, err=%v", err)
	}
}
