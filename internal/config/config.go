package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug bool `yaml:"is_debug"`

	Listen struct {
		BindIp string `yaml:"bind_ip"`
		Port   int    `yaml:"port"`
	}

	Database struct {
		Port         int    `yaml:"port"`
		Host         string `yaml:"host"`
		DatabaseName string `yaml:"database_name"`
		User         string `yaml:"user"`
		Password     string `yaml:"password"`
	} `yaml:"db"`
}

var (
	instance *Config
	once     sync.Once
)

func GetConfig(logger *log.Logger) Config {
	once.Do(func() {
		instance = &Config{}
		err := cleanenv.ReadConfig(ConfigFile, instance)
		if err != nil {
			logger.Fatalf("failed to read config, error: %v", err)
		}
	})
	return *instance
}
