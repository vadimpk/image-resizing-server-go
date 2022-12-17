package config

import (
	"github.com/spf13/viper"
	"strings"
	"time"
)

type (
	Config struct {
		Server ServerConfig
		Rabbit RabbitMQConfig
	}

	ServerConfig struct {
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
	}

	RabbitMQConfig struct {
		URL  string
		Port string
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
	if err := viper.UnmarshalKey("http", &cfg.Server); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("rabbitmq", &cfg.Rabbit); err != nil {
		return err
	}
	return nil
}
