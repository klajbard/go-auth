package config

import (
	"flag"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Services []ServicesConfig `yaml:"services"`
}

type ServicesConfig struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

var cfg = flag.String("cfg", "config.yaml", "config file path")
var Conf Configuration
var Channels = map[string]string{}

func (c *Configuration) GetConf() {
	flag.Parse()
	conf, err := ioutil.ReadFile(*cfg)
	if err != nil {
		log.Println(err)
	}

	if err := yaml.Unmarshal(conf, c); err != nil {
		log.Println(err)
	}
}
