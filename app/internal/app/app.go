package app

import (
	"cc/app/internal/config"
	"cc/app/internal/service"
	"cc/app/internal/storage"
	"cc/app/internal/transport"
	"cc/app/internal/transport/handler"
	"cc/app/pkg/postgres"
	"log"
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
	db, err := postgres.NewClient(postgres.Config{
		Host: app.cfg.Postgres.Host, Port: app.cfg.Postgres.Port, Database: app.cfg.Postgres.Database,
		Username: app.cfg.Postgres.Username, Password: app.cfg.Postgres.Password,
	})
	if err != nil {
		log.Fatal(err)
	}

	authStorage := storage.NewAuthStorage(db)
	authService := service.NewAuthService(authStorage, app.cfg.Auth.SigningKey, 1*time.Hour)

	statsStorage := storage.NewStatsStorage(db)
	statsService := service.NewStatsService(statsStorage)

	shortenStorage := storage.NewShortenStorage(db)
	shortenService := service.NewShortenService(shortenStorage, app.cfg.Shorten.Host)
	shortenHandler := handler.NewShortenHandler(shortenService, authService, statsService)

	userStorage := storage.NewUserStorage(db)
	userService := service.NewUserService(userStorage)
	userHandler := handler.NewUserHandler(userService, authService, shortenService)

	server := transport.New(
		shortenHandler,
		userHandler,
	)

	log.Fatal(server.Run(app.cfg.Server.Addr))
}
