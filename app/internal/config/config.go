package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

var (
	path = "config.yaml"
)

func New() (cfg Config) {
	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	return
}

type Config struct {
	Server     Server     `yaml:"server"`
	Postgres   Postgres   `yaml:"postgres"`
	Auth       Auth       `yaml:"auth"`
	Shorten    Shorten    `yaml:"shorten"`
	Prometheus Prometheus `yaml:"prometheus"`
	Redis      Redis      `yaml:"redis"`
}

type Server struct {
	Addr string `yaml:"addr" env:"HOST_ADDR" env-default:":8081"`
}

type Postgres struct {
	Host     string `yaml:"host" env:"POSTGRES_HOST" env-default:"localhost"`
	Port     string `yaml:"port" env:"POSTGRES_PORT" env-default:"5432"`
	Database string `yaml:"database" env:"POSTGRES_DATABASE" env-default:"postgres"`
	Username string `yaml:"username" env:"POSTGRES_USERNAME" env-default:"postgres"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-required:"true"`
}

type Auth struct {
	SigningKey string `yaml:"signing_key" env:"AUTH_SIGNING_KEY" env-required:"true"`
}

type Shorten struct {
	Host string `yaml:"host" env:"SHORTEN_HOST"`
}

type Prometheus struct {
	Addr string `yaml:"addr" env:"PROMETHEUS_ADDR" env-default:":8082"`
}

type Redis struct {
	Addr     string `yaml:"addr" env:"REDIS_ADDR" env-default:":6379"`
	Password string `yaml:"password" env:"REDIS_PASSWORD"`
}
