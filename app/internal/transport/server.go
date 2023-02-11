package transport

import (
	"cc/app/internal/transport/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

type Server struct {
	router   *gin.Engine
	handlers []Handler
}

func New(handlers ...Handler) *Server {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Use(middleware.Error())

	return &Server{handlers: handlers, router: router}
}

func (server *Server) Run(addr string) error {
	api := server.router.Group("/")

	for _, handler := range server.handlers {
		handler.Register(api)
	}

	return server.router.Run(addr)
}
