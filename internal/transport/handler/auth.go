package handler

import (
	"cc/internal/domain"
	"cc/internal/dto"
	service2 "cc/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthHandler struct {
	authService service2.AuthService
	userService service2.UserService
}

func NewAuthHandler(authService service2.AuthService, userService service2.UserService) *AuthHandler {
	return &AuthHandler{authService: authService, userService: userService}
}

func (handler *AuthHandler) Register(group *gin.RouterGroup) {
	group.POST("/signin", handler.SignIn)
	group.POST("/signup", handler.SignUp)
	group.POST("/refresh", handler.Refresh)
}

func (handler *AuthHandler) SignIn(c *gin.Context) {
	var request dto.Credentials
	if err := c.BindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}

	if err := request.Validate(); err != nil {
		_ = c.Error(err)
		return
	}

	user, err := handler.userService.SignIn(c,
		request,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var session domain.Session
	session, err = handler.authService.CreateSession(c,
		user.ID,
		c.ClientIP(),
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": session,
	})
}

func (handler *AuthHandler) SignUp(c *gin.Context) {
	var request dto.Credentials
	if err := c.BindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}

	if err := request.Validate(); err != nil {
		_ = c.Error(err)
		return
	}

	user, err := handler.userService.SignUp(c,
		request,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var session domain.Session
	session, err = handler.authService.CreateSession(c,
		user.ID,
		c.ClientIP(),
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": session,
	})
}

func (handler *AuthHandler) Refresh(c *gin.Context) {
	var request dto.Refresh
	if err := c.BindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}

	session, err := handler.authService.UpdateSession(c,
		request,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": session,
	})
}
