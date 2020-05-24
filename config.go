package main

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func (conf *Configuration) getConf() *Configuration {
	configFile := os.Getenv("CONF")
	if configFile == "" {
		configFile = "config.yml"
	}
	log.Printf("Application config file: %s", configFile)
	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return conf
}
