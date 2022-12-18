package rabbitmq

import (
	"github.com/spf13/viper"
	"strings"
)

type QueueConfig struct {
	Name          string `mapstructure:"name"`
	Durable       bool   `mapstructure:"durable"`
	DeleteUnused  bool   `mapstructure:"delete-unused"`
	Exclusive     bool   `mapstructure:"exclusive"`
	NoWait        bool   `mapstructure:"no-wait"`
	PrefetchCount int    `mapstructure:"prefetch-count"`
	AutoAck       bool   `mapstructure:"auto-ack"`
}

func Init(configPath string) (*QueueConfig, error) {
	if err := parseConfigPath(configPath); err != nil {
		return nil, err
	}

	var cfg QueueConfig
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

func unmarshal(cfg *QueueConfig) error {
	if err := viper.UnmarshalKey("rabbitmq-img-queue", &cfg); err != nil {
		return err
	}
	return nil
}
