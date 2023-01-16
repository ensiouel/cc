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
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"math/rand"
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

func GetRand[T any](s []T) T {
	return s[rand.Int()%len(s)]
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
		Addr:     app.cfg.Redis.Addr,
		Password: app.cfg.Redis.Password,
	})

	authStorage := storage.NewAuthStorage(client)
	authService := service.NewAuthService(authStorage, app.cfg.Auth.SigningKey, 1*time.Hour)

	statsStorage := storage.NewStatsStorage(client)
	statsService := service.NewStatsService(statsStorage)

	shorten := []uint64{}
	ua := []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/603.3.8 (KHTML, like Gecko) Version/10.1.2 Safari/603.3.8",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.90 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36 Edge/15.15063",
		"Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1; .NET CLR 1.1.4322) NS8/0.9.6",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.0 Mobile/14F89 Safari/602.1",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) FxiOS/8.1.1b4948 Mobile/14F89 Safari/603.2.4",
		"Mozilla/5.0 (iPad; CPU OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.0 Mobile/14F89 Safari/602.1",
		"Mozilla/5.0 (iPad; CPU OS 10_3_2 like Mac OS X) AppleWebKit/602.1.50 (KHTML, like Gecko) CriOS/58.0.3029.113 Mobile/14F89 Safari/602.1",
		"Opera/9.80 (Android; Opera Mini/28.0.2254/66.318; U; en) Presto/2.12.423 Version/12.16",
		"Mozilla/5.0 (Linux; Android 10; ONEPLUS A6003) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.0 Mobile Safari/537.36 EdgA/44.11.4.4140",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows Phone OS 7.0; Trident/3.1; IEMobile/7.0; NOKIA; Lumia 630)",
		"Mozilla/5.0 (compatible; SemrushBot/7~bl; +http://www.semrush.com/bot.html\", \"SemrushBot",
	}

	referer := []string{"https://vk.com", "https://youtube.com", "https://telegram.org", "https://google.com", "https://yandex.ru"}
	ip := []string{
		"197.153.142.126",
		"170.170.124.69",
		"52.92.126.13",
		"223.72.31.144",
		"168.150.22.220",
		"195.174.211.187",
		"251.27.90.166",
		"14.137.213.78",
		"73.232.249.220",
		"235.92.124.87",
		"178.103.65.225",
		"44.246.121.205",
		"186.49.136.31",
		"25.242.199.118",
		"189.241.34.162",
		"219.224.37.91",
		"120.150.53.253",
		"245.62.253.75",
		"89.158.29.0",
		"217.164.208.92",
	}

	{
		endDate := time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local)
		now := time.Now()

		for now.After(endDate) {
			endDate = endDate.Add(time.Duration(rand.Int63n(43200)) * time.Second)

			fmt.Println(endDate)

			err = statsService.CreateClickByUserAgent(ctx, endDate, GetRand(shorten), GetRand(ua), GetRand(referer), GetRand(ip))
			if err != nil {
				log.Fatal(err)
				return
			}
		}
	}

	shortenStorage := storage.NewShortenStorage(client)
	shortenService := service.NewShortenService(shortenStorage, app.cfg.Shorten.Host)

	userStorage := storage.NewUserStorage(client)
	userService := service.NewUserService(userStorage)

	authHandler := handler.NewAuthHandler(
		authService,
		userService,
	)

	shortenHandler := handler.NewShortenHandler(
		shortenService,
		authService,
		statsService,
		cache,
	)

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
		authHandler,
	)
	go func() {
		err = server.Run(app.cfg.Server.Addr)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Println(err)
		}
	}()

	<-ctx.Done()
}
