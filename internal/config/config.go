package config

import (
	"for9may/pkg/database"
	"for9may/pkg/redis"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type ServerCfg struct {
	HttpPort   uint16 `env:"SERVER_HTTP_PORT" env-default:"8080"`
	ServerMode string `env:"SERVER_MODE" env-default:"debug"`
}

type AdminCfg struct {
	Password string `env:"ADMIN_PASSWORD"`
	Login    string `env:"ADMIN_LOGIN"`
}

type Config struct {
	Server   ServerCfg       `env:"SERVER"`
	Admin    AdminCfg        `env:"ADMIN"`
	DataBase database.Config `env:"POSTGRES"`
	Redis    redis.Config    `env:"REDIS"`
}

func New() (*Config, error) {
	var cfg Config

	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
