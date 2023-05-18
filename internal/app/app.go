package app

import (
	"cc/internal/config"
	"cc/internal/service"
	"cc/internal/storage"
	"cc/internal/transport"
	"cc/internal/transport/handler"
	"cc/pkg/postgres"
	"context"
	"errors"
	"github.com/go-redis/redis/v9"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os/signal"
	"syscall"
)

type App struct {
	config config.Config
}

func New() *App {
	app := &App{
		config: config.New(),
	}

	return app
}

func (app *App) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pgClient, err := postgres.NewClient(ctx, postgres.Config{
		Host: app.config.Postgres.Host, Port: app.config.Postgres.Port, DB: app.config.Postgres.DB,
		User: app.config.Postgres.User, Password: app.config.Postgres.Password,
	})
	if err != nil {
		log.Fatal(err)
	}

	cache := redis.NewClient(&redis.Options{
		Addr: app.config.Redis.Addr,
	})
	ping := cache.Ping(ctx)
	if err = ping.Err(); err != nil {
		log.Fatal(err)
	}

	tagStorage := storage.NewTagStorage(pgClient)
	tagService := service.NewTagService(tagStorage)

	authStorage := storage.NewAuthStorage(pgClient)
	authService := service.NewAuthService(
		authStorage,
		app.config.Auth,
	)

	statsStorage := storage.NewStatsStorage(pgClient)
	statsService := service.NewStatsService(statsStorage)

	shortenStorage := storage.NewShortenStorage(pgClient)
	shortenService := service.NewShortenService(
		shortenStorage,
		app.config.Shorten.DomainURL,
	)

	userStorage := storage.NewUserStorage(pgClient)
	userService := service.NewUserService(userStorage)

	authHandler := handler.NewAuthHandler(
		authService,
		userService,
	)

	shortenHandler := handler.NewShortenHandler(
		shortenService,
		authService,
		tagService,
		statsService,
	)

	userHandler := handler.NewUserHandler(
		userService,
		authService,
		shortenService,
		tagService,
	)

	redirectHandler := handler.NewRedirectHandler(
		shortenService,
		statsService,
		cache,
		app.config.Shorten.DefaultURL,
	)

	go func() {
		http.Handle("/metrics", promhttp.Handler())

		err = http.ListenAndServe(app.config.Prometheus.Addr, nil)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	go func() {
		err = transport.New().
			Handle(
				shortenHandler,
				userHandler,
				authHandler,
				redirectHandler,
				authService,
			).
			Run(app.config.Server.Addr)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
}
