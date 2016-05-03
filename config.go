package main

import (
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var cfg Config

type Config struct {
	SMTP
	Log
}
type SMTP struct {
	Port  int
	Hosts []string
	Accept
}
type Accept struct {
	Domains []string
}
type Log struct {
	Path string
}

func prepareConfig() {
	cfgRaw, err := ioutil.ReadFile("config.yaml")
	if err == nil {
		yaml.Unmarshal(cfgRaw, &cfg)
	}
	defaultCfg := Config{
		SMTP: SMTP{
			Port:   25,
			Accept: Accept{Domains: []string{"localhost"}},
		},
		Log: Log{
			Path: "smtpd.log",
		},
	}
	err = mergo.Merge(&cfg, defaultCfg)
	if err != nil {
		log.Fatal(err)
	}
}
