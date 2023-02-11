package handler

import (
	"cc/app/internal/domain"
	"cc/app/internal/dto"
	"cc/app/internal/service"
	"cc/app/internal/transport"
	"cc/app/internal/transport/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type UserHandler struct {
	userService    service.UserService
	authService    service.AuthService
	shortenService service.ShortenService
	tagService     service.TagService
}

func NewUserHandler(userService service.UserService, authService service.AuthService, shortenService service.ShortenService, tagService service.TagService) transport.Handler {
	return &UserHandler{userService: userService, authService: authService, shortenService: shortenService, tagService: tagService}
}

func (handler *UserHandler) Register(group *gin.RouterGroup) {
	authorized := group.Group("/")
	authorized.Use(middleware.Auth(handler.authService))
	{
		authorized.GET("/api/users/:id", handler.GetUser)
		authorized.GET("/api/users/:id/shortens", handler.SelectUserShortens)

		authorized.GET("/api/users/:id/tags", handler.SelectUserTags)
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
	var request dto.SelectShortens
	if err := c.BindQuery(&request); err != nil {
		_ = c.Error(err)
		return
	}

	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	var shortens []domain.Shorten
	if len(request.Tags) > 0 {
		shortens, err = handler.shortenService.SelectByTags(c,
			userID,
			request.Tags,
		)
		if err != nil {
			_ = c.Error(err)
			return
		}
	} else {
		shortens, err = handler.shortenService.SelectByUser(c,
			userID,
		)
		if err != nil {
			_ = c.Error(err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"response": shortens,
	})
}

func (handler *UserHandler) SelectUserTags(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	var tags domain.Tags
	tags, err = handler.tagService.SelectByUser(c, userID)

	c.JSON(http.StatusOK, gin.H{
		"response": tags,
	})
}
