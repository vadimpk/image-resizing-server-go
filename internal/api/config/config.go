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
		Main
		FileStorage
	}

	ServerConfig struct {
		Port          string        `mapstructure:"port"`
		ReadTimeout   time.Duration `mapstructure:"readTimeout"`
		WriteTimeout  time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMB   int           `mapstructure:"maxHeaderBytes"`
		MaxFileSizeMB int64         `mapstructure:"maxFileSizeMB"`
	}

	RabbitMQConfig struct {
		URL  string
		Port string
	}

	FileStorage struct {
		DirPath string `mapstructure:"dir"`
	}

	Main struct {
		Timeout time.Duration `mapstructure:"graceful-timeout"`
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
	if err := viper.UnmarshalKey("main", &cfg.Main); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("file-storage", &cfg.FileStorage); err != nil {
		return err
	}
	return nil
}
