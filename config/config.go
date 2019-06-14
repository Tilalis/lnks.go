package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Config initial support
type Config struct {
	Database  string `json:"database"`
	SecretKey string `json:"secretkey"`
}

// LoadConfig function to load config
func LoadConfig(filename string) (*Config, error) {
	jsonFile, err := os.Open(filename)

	if err != nil {
		return nil, ErrConfigFileNotFound.SetFile(filename)
	}

	defer jsonFile.Close()

	jsonBytes, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		return nil, ErrReadingConfigFile.SetFile(filename)
	}

	var config Config
	err = json.Unmarshal(jsonBytes, &config)

	if err != nil {
		return nil, ErrMalformedConfigFile.SetFile(filename)
	}

	return &config, nil
}
