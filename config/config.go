package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type (
	HTTP struct {
		Address           string        `yaml:"address"`
		BasePath          string        `yaml:"base_path"`
		ReadHeaderTimeout time.Duration `yaml:"read_header_timeout"`
		ReadTimeout       time.Duration `yaml:"read_timeout"`
		WriteTimeout      time.Duration `yaml:"write_timeout"`
		IdleTimeout       time.Duration `yaml:"idle_timeout"`
		ShutdownTimeout   time.Duration `yaml:"shutdown_timeout"`
		AllowedOrigins    []string      `yaml:"allowed_origins"`
	}

	DB struct {
		Kind           string `yaml:"kind"`
		Host           string `yaml:"host"`
		Port           string `yaml:"port"`
		UserEnvKey     string `yaml:"user_env_key"`
		PassEnvKey     string `yaml:"pass_env_key"`
		DBName         string `yaml:"db_name"`
		MaxConnections int    `yaml:"max_connections"`
		Timeout        int    `yaml:"timeout"`
	}

	Redis struct {
		Address  string `yaml:"address"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	}

	Logging struct {
		Level string `yaml:"level"`
	}

	Config struct {
		HTTP    HTTP    `yaml:"http"`
		DB      DB      `yaml:"db"`
		Redis   Redis   `yaml:"redis"`
		Logging Logging `yaml:"logging"`
	}
)

func Read(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
