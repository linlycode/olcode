package main

import yaml "gopkg.in/yaml.v2"

type config struct {
	Port   int    `yaml:"port"`
	DBPath string `yaml:"db"`
}

func loadConfig(data []byte) (*config, error) {
	c := config{}
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}
