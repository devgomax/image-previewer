package config

import (
	"net"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Key строковый алиас для ключей конфигурации.
type Key = string

const (
	// ConfigPath ключ конфигурации, указывающий путь до конфиг-файла.
	ConfigPath Key = "config"
)

// LoggerConfig модель конфига для логгера.
type LoggerConfig struct {
	Level              string `mapstructure:"level"`
	DisableSampling    bool   `mapstructure:"disable_sampling"`
	TimestampFieldName string `mapstructure:"timestamp_field_name"`
	LevelFieldName     string `mapstructure:"level_field_name"`
	MessageFieldName   string `mapstructure:"message_field_name"`
	ErrorFieldName     string `mapstructure:"error_field_name"`
	TimeFieldFormat    string `mapstructure:"time_field_format"`
}

// ServerConfig модель конфига для сервера.
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

// GetAddr возвращает строку вида "host:port".
func (sc *ServerConfig) GetAddr() string {
	return net.JoinHostPort(sc.Host, sc.Port)
}

// LRUCacheConfig модель конфига для кэша LRU.
type LRUCacheConfig struct {
	Size int `mapstructure:"size"`
}

// Config модель основного конфига приложения.
type Config struct {
	Logger         LoggerConfig   `mapstructure:"logger"`
	HTTPConfig     ServerConfig   `mapstructure:"http"`
	LRUCacheConfig LRUCacheConfig `mapstructure:"lru_cache"`
}

// NewConfig конструктор для основного конфига приложения.
func NewConfig() (*Config, error) {
	configPath := pflag.String(ConfigPath, "/etc/image_previewer/config.toml", "Path to configuration file")
	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return nil, errors.Wrap(err, "[config::NewConfig]: failed to bind flag set to config")
	}

	var c Config

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
	viper.SetConfigFile(*configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "[config::NewConfig]: failed to discover and read config file")
	}

	err := viper.Unmarshal(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
