package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

var (
	Port    = "3000"
	Release = "development"
	User    = "Anon"
	Time    = time.Now().Format(time.RFC3339)
	Hash    = "N/A"
	Module  = fx.Options(fx.Provide(NewConfig))
)

type (
	HTTP struct {
		Port string
	}

	Config struct {
		HTTP      `yaml:"http"`
		Hash      string
		Time      string
		User      string
		Release   string
	}
)

func NewConfig() *Config {

	cfg := &Config{
		HTTP: HTTP{
			Port: Port,
		},
		Hash:      Hash,
		Time:      Time,
		User:      User,
		Release:   Release,
	}

	return cfg
}

func (cfg *Config) InitConfig(dir string) error {
	err := godotenv.Load()

	if err != nil {
		return err
	}

	port, isFound := os.LookupEnv("PORT")

	if isFound {
		cfg.HTTP.Port = port
	}

  release, isFound := os.LookupEnv("RELEASE")

  if isFound {
    cfg.Release = release
  }

	return nil
}
