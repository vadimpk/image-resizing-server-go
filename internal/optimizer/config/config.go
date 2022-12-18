package config

import (
	"github.com/spf13/viper"
	"strings"
	"time"
)

type (
	Config struct {
		Rabbit RabbitMQConfig
		Main
		FileStorage
	}

	RabbitMQConfig struct {
		URL  string
		Port string
	}

	Main struct {
		Timeout time.Duration `mapstructure:"graceful-timeout"`
	}

	FileStorage struct {
		DirPath string `mapstructure:"dir"`
	}
)

func Init(configPath string) (*Config, error) {
	if err := parseConfigPath(configPath); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func parseConfigPath(filepath string) error {
	path := strings.Split(filepath, "/")

	viper.AddConfigPath(path[0])
	viper.SetConfigName(path[1])

	return viper.ReadInConfig()
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("rabbitmq", &cfg.Rabbit); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("main", &cfg.Main); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("file-storage", &cfg.FileStorage); err != nil {
		return err
	}
	return nil
}
