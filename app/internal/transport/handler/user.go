package handler

import (
	"cc/app/internal/domain"
	"cc/app/internal/service"
	"cc/app/internal/transport"
	"cc/app/internal/transport/middleware/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type UserHandler struct {
	userService    service.UserService
	authService    service.AuthService
	shortenService service.ShortenService
}

func NewUserHandler(userService service.UserService, authService service.AuthService, shortenService service.ShortenService) transport.Handler {
	return &UserHandler{userService: userService, authService: authService, shortenService: shortenService}
}

func (handler *UserHandler) Register(group *gin.RouterGroup) {
	authorized := group.Group("/")
	authorized.Use(auth.Middleware(handler.authService))
	{
		authorized.GET("/api/users/:id", handler.GetUser)
		authorized.GET("/api/users/:id/shortens", handler.SelectUserShortens)
	}
}

func (handler *UserHandler) GetUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	var user domain.User
	user, err = handler.userService.GetUserByID(c,
		userID,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": user,
	})
}

func (handler *UserHandler) SelectUserShortens(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	var shortens []domain.Shorten
	shortens, err = handler.shortenService.SelectShortensByUserID(c,
		userID,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": shortens,
	})
}
