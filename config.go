package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type config struct {
	cyberpower []configcyberpower
}

type configcyberpower struct {
	host     string
	username string
	password string
}

func read_config(filename string) *config {

	c := &config{}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("unable to parse %s: %s", filename, err)
		return c
	}

	err = yaml.Unmarshal(data, c)
	if err != nil {
		log.Printf("unable to parse %s: %s", filename, err)
	}

	return c
}
