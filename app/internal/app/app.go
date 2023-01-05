package app

import (
	"cc/app/internal/config"
	"cc/app/internal/service"
	"cc/app/internal/storage"
	"cc/app/internal/transport"
	"cc/app/internal/transport/handler"
	"cc/app/pkg/postgres"
	"context"
	"errors"
	"github.com/go-redis/redis/v9"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	cfg config.Config
}

func New() (app *App) {
	app = &App{}
	app.cfg = config.New()
	return
}

func (app *App) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	client, err := postgres.NewClient(ctx, postgres.Config{
		Host: app.cfg.Postgres.Host, Port: app.cfg.Postgres.Port, Database: app.cfg.Postgres.Database,
		Username: app.cfg.Postgres.Username, Password: app.cfg.Postgres.Password,
	})
	if err != nil {
		log.Fatal(err)
	}

	cache := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	authStorage := storage.NewAuthStorage(client)
	authService := service.NewAuthService(authStorage, app.cfg.Auth.SigningKey, 1*time.Hour)

	statsStorage := storage.NewStatsStorage(client)
	statsService := service.NewStatsService(statsStorage)

	shortenStorage := storage.NewShortenStorage(client)
	shortenService := service.NewShortenService(shortenStorage, app.cfg.Shorten.Host)
	shortenHandler := handler.NewShortenHandler(
		shortenService,
		authService,
		statsService,
		cache,
	)

	userStorage := storage.NewUserStorage(client)
	userService := service.NewUserService(userStorage)
	userHandler := handler.NewUserHandler(
		userService,
		authService,
		shortenService,
	)

	go func() {
		err = http.ListenAndServe(app.cfg.Prometheus.Addr, promhttp.Handler())
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Println(err)
		}
	}()

	server := transport.New(
		shortenHandler,
		userHandler,
	)
	go func() {
		err = server.Run(app.cfg.Server.Addr)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Println(err)
		}
	}()

	<-ctx.Done()
}
