package main

import (
	"errors"
	"github.com/aryann/difflib"
	"github.com/imdario/mergo"
	"github.com/nochso/mlog"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strings"
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
	err = mergo.Merge(&cfg, getDefaultConfig())
	if err != nil {
		log.Fatal(errors.New("unable to merge configuration: " + err.Error()))
	}
}

func getDefaultConfig() *Config {
	return &Config{
		SMTP: SMTP{
			Port:   25,
			Accept: Accept{Domains: []string{"localhost"}},
		},
		Log: Log{
			Path: "smtpd.log",
		},
	}
}

func getConfigDiff() string {
	configYaml, err := yaml.Marshal(&cfg)
	if err != nil {
		mlog.Error(errors.New("Unable to dump configuration to YAML: " + err.Error()))
		return ""
	}
	defaultYaml, err := yaml.Marshal(getDefaultConfig())
	diff := difflib.Diff(strings.Split(string(defaultYaml), "\n"), strings.Split(string(configYaml), "\n"))
	lines := []string{}
	for _, d := range diff {
		if d.Delta == difflib.LeftOnly {
			continue
		}
		lines = append(lines, d.String())
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}
