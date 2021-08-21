package config

import (
	conf "github.com/spf13/viper"
)

type Config struct {
	Collectors struct {
		MaxTimeStorage int `yaml:"maxTimeStorage"`
	} `yaml:"collector"`
	Log struct {
		Level string `yaml:"level"`
	} `yaml:"log"`
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	Metric map[string]bool `yaml:"metric"`
}

func LoadCfg() (Config, error) {
	config := Config{}
	conf.SetConfigName("config")
	conf.AddConfigPath("./cfg/")
	if err := conf.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := conf.Unmarshal(&config); err != nil {
		panic(err)
	}

	return config, nil
}
