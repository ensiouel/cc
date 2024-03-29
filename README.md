## Deployment

**Build** application

```shell
docker compose build
```

**Run** application

```shell
docker compose up -d
```

## Config

```dotenv
GIN_MODE=release

SERVER_ADDR=:8080
PROMETHEUS_ADDR=:9091

POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=cc

REDIS_ADDR=redis:6379

AUTH_EXPIRATION_AT=1h
AUTH_SIGNING_KEY=SIGNING_KEY

SHORTEN_DOMAIN_URL=localhost:8081
SHORTEN_DEFAULT_URL=https://www.google.com
```