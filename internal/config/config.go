package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
    DbHost     string `env:"DB_HOST" envDefault:"localhost"`
    DbPort     string `env:"DB_PORT" envDefault:"5432"`
    DbUser     string `env:"DB_USER,required"`
    DbPwd      string `env:"DB_PWD,required"`
    DbName     string `env:"DB_NAME,required"`
    JwtSecret  string `env:"JWT_SECRET,required"`
    SessionKey string `env:"SESSION_KEY,required"`
    ServerHost string `env:"SERVER_HOST" envDefault:"0.0.0.0"`
    ServerPort string `env:"SERVER_PORT" envDefault:"8080"`
}

func LoadConfig() (*Config, error) {
  cfg := &Config{}
  err := env.Parse(cfg)
	if err != nil {
    return nil, err
  }
  return cfg, nil
}
