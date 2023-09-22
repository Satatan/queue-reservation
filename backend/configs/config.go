package configs

import (
	"queue_reservation/pkg/logx"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App appConfigs
}

type appConfigs struct {
	Port string `envconfig:"APP_PORT"`
}

var config Config

func Init() *Config {
	err := godotenv.Load()

	if err != nil {
		envFileNotFound := strings.Contains(err.Error(), "no such file or directory")
		if !envFileNotFound {
			logx.GetLog().Fatalf("read config err := %v", err)
		} else {
			logx.GetLog().Info("use environment from OS")
		}
	}
	if err = envconfig.Process("", &config); err != nil {
		logx.GetLog().Fatalf("parse configs error: %v", err)
	}
	return &config

}
