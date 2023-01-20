package handler

import (
	"cc/app/internal/apperror"
	"cc/app/internal/domain"
	"cc/app/internal/dto"
	"cc/app/internal/service"
	"cc/app/internal/transport"
	"cc/app/internal/transport/middleware/auth"
	"cc/app/pkg/base62"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type ShortenHandler struct {
	shortenService service.ShortenService
	authService    service.AuthService
	statsService   service.StatsService
	cache          *redis.Client
}

func NewShortenHandler(shortenService service.ShortenService, authService service.AuthService, statsService service.StatsService, cache *redis.Client) transport.Handler {
	return &ShortenHandler{shortenService: shortenService, authService: authService, statsService: statsService, cache: cache}
}

func (handler *ShortenHandler) Register(group *gin.RouterGroup) {
	group.GET("/:key", handler.Redirect)

	authorized := group.Group("/")
	authorized.Use(auth.Middleware(handler.authService))
	{
		authorized.GET("/api/shortens/:key/stats", handler.GetShortenStats)

		authorized.POST("/api/shortens", handler.CreateShorten)
		authorized.GET("/api/shortens/:key", handler.GetShorten)
		authorized.PATCH("/api/shortens/:key", handler.UpdateShorten)
		authorized.DELETE("/api/shortens/:key", handler.DeleteShorten)
	}
}

func (handler *ShortenHandler) Redirect(c *gin.Context) {
	shortenKey := c.Param("key")
	shortenID, err := base62.Decode(shortenKey)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	var url string
	url, err = handler.cache.Get(c, "shorten:"+shortenKey).Result()
	if err != nil && errors.Is(err, redis.Nil) == false {
		_ = c.Error(apperror.ErrInternalError.SetError(err))
		return
	}

	if url == "" {
		url, err = handler.shortenService.GetShortenURL(c, shortenID)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		handler.cache.Set(c, "shorten:"+shortenKey, url, 1*time.Hour)
	}

	now := time.Now()

	err = handler.statsService.CreateClickByUserAgent(c,
		now,
		shortenID,
		c.Request.Header.Get("User-Agent"),
		c.Request.Referer(),
		c.ClientIP(),
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Redirect(http.StatusSeeOther, url)
}

func (handler *ShortenHandler) GetShorten(c *gin.Context) {
	shortenKey := c.Param("key")
	shortenID, err := base62.Decode(shortenKey)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var shorten domain.Shorten
	shorten, err = handler.shortenService.GetShortenByID(c,
		shortenID,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": shorten,
	})
}

func (handler *ShortenHandler) CreateShorten(c *gin.Context) {
	var request dto.CreateShorten
	if err := c.BindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}

	if err := request.Validate(); err != nil {
		_ = c.Error(err)
		return
	}

	var userID uuid.UUID
	if value, exists := c.Get("user_id"); exists && value != nil {
		userID = value.(uuid.UUID)
	}

	shorten, err := handler.shortenService.CreateShorten(c,
		userID,
		request,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": shorten,
	})
}

func (handler *ShortenHandler) UpdateShorten(c *gin.Context) {
	var request dto.UpdateShorten
	if err := c.BindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}

	if err := request.Validate(); err != nil {
		_ = c.Error(err)
		return
	}

	var userID uuid.UUID
	if value, exists := c.Get("user_id"); exists && value != nil {
		userID = value.(uuid.UUID)
	}

	shortenKey := c.Param("key")
	shortenID, err := base62.Decode(shortenKey)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var shorten domain.Shorten
	shorten, err = handler.shortenService.UpdateShorten(c,
		userID,
		shortenID,
		request,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": shorten,
	})
}

func (handler *ShortenHandler) DeleteShorten(c *gin.Context) {
	shortenID, err := base62.Decode(c.Param("key"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	var userID uuid.UUID
	if value, exists := c.Get("user_id"); exists && value != nil {
		userID = value.(uuid.UUID)
	}

	err = handler.shortenService.DeleteShorten(c,
		userID,
		shortenID,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": 1,
	})
}

func (handler *ShortenHandler) GetShortenStats(c *gin.Context) {
	var request dto.GetShortenStats
	if err := c.BindQuery(&request); err != nil {
		_ = c.Error(err)
		return
	}

	if err := request.Validate(); err != nil {
		_ = c.Error(err)
		return
	}

	shortenKey := c.Param("key")
	shortenID, err := base62.Decode(shortenKey)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var stats domain.Stats
	stats, err = handler.statsService.GetStats(c,
		shortenID,
		request,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": stats,
	})
}
