package main

import (
	"os"

	"github.com/ghodss/yaml"
)

// Config : configuration struct
type Config struct {
	Sqldb string `json:"sqldb"`
}

func config() (sqldb string, err error) {
	var v Config
	configFile, err := os.ReadFile("/etc/soh-router/config.yaml")
	check(err)
	err = yaml.Unmarshal(configFile, &v)
	sqldb = v.Sqldb
	return sqldb, err
}
