package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Db struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
	}
	Server struct {
		Port int
	}
}

func NewConfig() Config {
	wd, err := os.Getwd()

	confData, err := os.ReadFile(fmt.Sprintf("%s/./.env/config.yaml", wd))

	if err != nil {
		log.Fatal(err)
	}

	var config Config
	err = yaml.Unmarshal(confData, &config)

	if err != nil {
		log.Fatal(err)
	}

	return config
}
