package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"time"
)

type Config struct {
	Server     Server
	Prometheus Prometheus
	Postgres   Postgres
	Auth       Auth
	Redis      Redis
	Shorten    Shorten
}

type Server struct {
	Addr string `env:"SERVER_ADDR"`
}

type Prometheus struct {
	Addr string `env:"PROMETHEUS_ADDR"`
}

type Postgres struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	DB       string `env:"POSTGRES_DB"`
}

type Auth struct {
	ExpirationAt time.Duration `env:"AUTH_EXPIRATION_AT"`
	SigningKey   string        `env:"AUTH_SIGNING_KEY" env-required:"true"`
}

type Redis struct {
	Addr string `env:"REDIS_ADDR"`
}

type Shorten struct {
	DomainURL  string `env:"SHORTEN_DOMAIN_URL"`
	DefaultURL string `env:"SHORTEN_DEFAULT_URL"`
}

func New() Config {
	var config Config
	err := cleanenv.ReadEnv(&config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}
