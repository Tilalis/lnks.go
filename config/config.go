package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Config initial support
type Config struct {
	Server struct {
		Host string `json:"host"`
		Port string `json:"port"`
	} `json:"server"`
	Address  string `json:"-"`
	Database struct {
		Driver           string `json:"driver"`
		ConnectionString string `json:"connectionString"`
	} `json:"database"`
	SecretKey string `json:"secretKey"`
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

	config.Address = fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port)

	return &config, nil
}
