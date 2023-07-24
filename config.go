package main

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
)

// Config : configuration struct
type Config struct {
	Sqldb string `json:"sqldb"`
}

func config() (sqldb string) {
	var v Config
	configFile, err := ioutil.ReadFile("/etc/soh-router/config.yaml")
	check(err)
	yaml.Unmarshal(configFile, &v)
	sqldb = v.Sqldb
	return
}
