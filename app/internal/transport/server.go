package transport

import (
	"cc/app/internal/transport/middleware/errs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"time"
)

type Server struct {
	router   *gin.Engine
	handlers []Handler
}

func New(handlers ...Handler) *Server {
	router := gin.Default()

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Use(errs.Middleware())

	return &Server{handlers: handlers, router: router}
}

func (server *Server) Run(addr string) error {
	api := server.router.Group("/")

	for _, handler := range server.handlers {
		handler.Register(api)
	}

	return server.router.Run(addr)
}
