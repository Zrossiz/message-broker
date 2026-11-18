package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Zrossiz/message-broker/internal/helpers"
	"gopkg.in/yaml.v3"
)

type Config struct {
	App  App  `yaml:"app"`
	Repo Repo `yaml:"repo"`
}

type App struct {
	LoggerLevel string `yaml:"logger_level"`
	Addr        string `yaml:"address"`
	Port        int    `yaml:"port"`
}

type Repo struct {
	Filepath string `yaml:"filepath"`
}

func New() (*Config, error) {
	configsPath, err := helpers.WithBasePath("configs")
	if err != nil {
		return nil, fmt.Errorf("base path for public: %w", err)
	}

	env := os.Getenv("GOLANG_ENV")
	if env == "" {
		env = "dev"
	} else {
		env = strings.ToLower(env)
	}

	filename := fmt.Sprintf("%s/%s.yml", configsPath, env)
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		filename = configsPath + "/dev.yaml"
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
