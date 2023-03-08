package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Url       string `json:"url"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	ClassName string `json:"class_name"`
}

func Load(filepath string) Config {
	var config Config

	configFile, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}

	defer configFile.Close()

	if err := json.NewDecoder(configFile).Decode(&config); err != nil {
		panic(err)
	}

	return config
}
