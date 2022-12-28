package handler

import (
	"cc/app/internal/domain"
	"cc/app/internal/dto"
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
	group.POST("/api/users/signin", handler.SignIn)
	group.POST("/api/users/signup", handler.SignUp)
	group.POST("/api/users/refresh", handler.Refresh)

	authorized := group.Group("/")
	authorized.Use(auth.Middleware(handler.authService))
	{
		authorized.GET("/api/users/:id", handler.GetUser)
		authorized.GET("/api/users/:id/shortens", handler.SelectUserShortens)
	}
}

func (handler *UserHandler) SignIn(c *gin.Context) {
	var request dto.Credentials

	if err := c.BindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}

	if err := request.Validate(); err != nil {
		_ = c.Error(err)
		return
	}

	user, err := handler.userService.SignIn(c, request)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var session domain.Session
	session, err = handler.authService.CreateSession(c, user.ID, c.ClientIP())
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.SetCookie("refresh_token", session.RefreshToken.String(), 86400,
		"/", c.Request.Host, false, true)

	c.JSON(http.StatusOK, gin.H{
		"response": session,
	})
}

func (handler *UserHandler) SignUp(c *gin.Context) {
	var request dto.Credentials

	if err := c.BindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}

	if err := request.Validate(); err != nil {
		_ = c.Error(err)
		return
	}

	user, err := handler.userService.SignUp(c, request)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var session domain.Session
	session, err = handler.authService.CreateSession(c, user.ID, c.ClientIP())
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.SetCookie("refresh_token", session.RefreshToken.String(), 86400,
		"/", c.Request.Host, false, true)

	c.JSON(http.StatusOK, gin.H{
		"response": session,
	})
}

func (handler *UserHandler) Refresh(c *gin.Context) {
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		_ = c.Error(err)
		return
	}

	var refreshToken uuid.UUID
	refreshToken, err = uuid.Parse(cookie)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var session domain.Session
	session, err = handler.authService.UpdateSession(c, refreshToken)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.SetCookie("refresh_token", session.RefreshToken.String(), 86400,
		"/", c.Request.Host, false, true)

	c.JSON(http.StatusOK, gin.H{
		"response": session,
	})
}

func (handler *UserHandler) GetUser(c *gin.Context) {
	var (
		userID uuid.UUID
		err    error
	)

	if userID, err = uuid.Parse(c.Param("id")); err != nil {
		_ = c.Error(err)
		return
	}

	var user domain.User
	user, err = handler.userService.GetUserByID(c, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": user,
	})
}

func (handler *UserHandler) SelectUserShortens(c *gin.Context) {
	var (
		userID uuid.UUID
		err    error
	)

	if userID, err = uuid.Parse(c.Param("id")); err != nil {
		_ = c.Error(err)
		return
	}

	var shortens []domain.Shorten
	shortens, err = handler.shortenService.SelectShortensByUserID(c, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": shortens,
	})
}
