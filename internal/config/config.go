package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	URL      string  `json:"url"`
	Selector string  `json:"selector"`
	Output   string  `json:"output"`
}

func LoadConfig(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return Config{}, err
	}
	return config, nil
}
