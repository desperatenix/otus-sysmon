package config

import (
	conf "github.com/spf13/viper"
	"log"
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

func LoadCfg(cfgPath string) (Config, error) {
	config := Config{}
	//conf.SetConfigName(cfgName)
	conf.AddConfigPath("/etc/sysmon/")
	conf.AddConfigPath("./cfg/")
	conf.AddConfigPath(cfgPath)



	if err := conf.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	if err := conf.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}

	return config, nil
}
