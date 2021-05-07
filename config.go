package main

import (
	"log"
	"io/ioutil"
	"regexp"
	"gopkg.in/yaml.v2"
)

type vfioConfig struct {
	Name string      `yaml:"resourceName"`
	Vendor string    `yaml:"vendorId"`
	Device []string  `yaml:"deviceId"`
}

func readConfigFile(fileName string) []vfioConfig {
	var config []vfioConfig

	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatal(err)
	}

	re := regexp.MustCompile(`^[a-z0-9.-]+\/[a-z0-9.-]+$`)
	for _,group := range config {
		if re.MatchString(group.Name) == false {
			log.Fatal("Invalid resourceName " + group.Name)
		}
	}

	return config
}
