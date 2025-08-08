package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Github GithubConfig `mapstructure:"github"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type GithubConfig struct {
	CommitURLList []string `mapstructure:"commit-url-list"`
	Authorization string   `mapstructure:"authorization"`
	Accept        string   `mapstructure:"accept"`
	APIVersion    string   `mapstructure:"api-version"`
}

var AppConfig Config

func LoadConfig(configPath string) error {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		return err
	}

	return nil
}
