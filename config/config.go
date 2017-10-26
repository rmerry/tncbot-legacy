// Package config is used for parsing a tncbot configuration file
package config

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

// Config defines the configuration options for tncbot
type Config struct {
	Password string `json:"password"`
	Channel  string `json:"channel"`
	Ident    string `json:"ident"`
	Nickname string `json:"nickname"`
	Port     int    `json:"port"`
	Server   string `json:"server"`
}

// Parse reads a configuration file and returns
func Parse() (*Config, error) {
	configFile, err := filepath.Abs("config.json")
	if err != nil {

	}

	configBytes, err := ioutil.ReadFile(configFile)
	if err != nil {

	}

	var c Config
	err = json.Unmarshal(configBytes, &c)
	if err != nil {

	}

	return &c, nil
}
