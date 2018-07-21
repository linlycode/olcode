package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/linlycode/olcode/pkg/api"
	"github.com/linlycode/olcode/pkg/common"
)

var configFilePath = flag.String("config", "", "config file")

func main() {
	flag.Parse()
	configData, err := ioutil.ReadFile(*configFilePath)
	common.Assertf(err == nil, "cannot open config file %s", *configFilePath)
	c, err := loadConfig(configData)
	common.Assertf(err == nil, "fail to load config, err=%v", err)

	log.Printf("serve on port %d", c.Port)
	s := api.NewService(c.Port)
	if err := s.Serve(); err != nil {
		log.Printf("server stopped, err=%v", err)
	}
}
