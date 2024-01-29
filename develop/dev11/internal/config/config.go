package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type (
	HTTPConfig struct {
		Port int `json:"port"`
	}

	Config struct {
		HTTP HTTPConfig `json:"http"`
	}
)

func Read(fileName string) (*Config, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("read config filer error: %w", err)
	}

	c := Config{}
	err = json.Unmarshal(data, &c)
	if err != nil {
		return nil, fmt.Errorf("unmarshal config file error: %w", err)
	}

	return &c, nil
}
