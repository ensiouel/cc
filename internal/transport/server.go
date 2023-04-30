package transport

import (
	"cc/internal/service"
	"cc/internal/transport/handler"
	"cc/internal/transport/middleware"
	"github.com/chenjiandongx/ginprom"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

func New() *Server {
	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.ExposeHeaders = []string{"*"}

	router.Use(
		cors.New(corsConfig),
		ginprom.PromMiddleware(&ginprom.PromOpts{
			EndpointLabelMappingFn: func(c *gin.Context) string {
				return c.FullPath()
			},
		}),
	)

	return &Server{router: router}
}

func (server *Server) Handle(
	shortenHandler *handler.ShortenHandler,
	userHandler *handler.UserHandler,
	authHandler *handler.AuthHandler,
	redirectHandler *handler.RedirectHandler,
	authService service.AuthService,
) *Server {
	redirectHandler.Register(server.router.Group("/"))

	api := server.router.Group("/api", middleware.Error())
	{
		authHandler.Register(api.Group("/auth"))

		authorized := api.Group("/", middleware.Auth(authService))
		{
			shortenHandler.Register(authorized.Group("/shortens"))
			userHandler.Register(authorized.Group("/users"))
		}

	}

	return server
}

func (server *Server) Run(addr string) error {
	return server.router.Run(addr)
}
