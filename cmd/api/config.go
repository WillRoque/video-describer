package main

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server struct {
		Host  string `env:"SRV_HOST" env-description:"Server host" env-default:"0.0.0.0"`
		Port  int    `env:"SRV_PORT" env-description:"Server port" env-default:"8080"`
		Debug bool   `env:"SRV_DEBUG" env-description:"Enable debug logs" env-default:"false"`
	}
	OpenAI struct {
		ApiKey string `env:"OPENAI_API_KEY" env-description:"OpenAI's API key" env-required:"true"`
	}
}

func readConfig() Config {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		msg := "Configuration missing"
		envHelp, _ := cleanenv.GetDescription(&cfg, &msg)
		fmt.Println(envHelp)
		os.Exit(1)
	}

	return cfg
}
