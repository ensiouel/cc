package transport

import (
	"cc/app/internal/transport/middleware/errs"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router   *gin.Engine
	handlers []Handler
}

func New(handlers ...Handler) *Server {
	router := gin.Default()

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
