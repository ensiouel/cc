package handler

import (
	"cc/internal/service"
	"cc/pkg/apperror"
	"cc/pkg/base62"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"time"
)

type RedirectHandler struct {
	shortenService service.ShortenService
	statsService   service.StatsService
	cache          *redis.Client
	defaultURL     string
}

func NewRedirectHandler(
	shortenService service.ShortenService,
	statsService service.StatsService,
	cache *redis.Client,
	defaultURL string,
) *RedirectHandler {
	return &RedirectHandler{
		shortenService: shortenService,
		statsService:   statsService,
		cache:          cache,
		defaultURL:     defaultURL,
	}
}

func (handler *RedirectHandler) Register(group *gin.RouterGroup) {
	group.GET("/:key", handler.Redirect)
}

func (handler *RedirectHandler) Redirect(c *gin.Context) {
	shortenKey := c.Param("key")
	shortenID, err := base62.Decode(shortenKey)
	if err != nil {
		log.Println(err)
		c.Redirect(http.StatusSeeOther, handler.defaultURL)
		return
	}

	var url string
	url, err = handler.cache.Get(c, "shorten:"+shortenKey).Result()
	if err != nil && errors.Is(err, redis.Nil) == false {
		log.Println(err)
		c.Redirect(http.StatusSeeOther, handler.defaultURL)
		return
	}

	if url == "" {
		url, err = handler.shortenService.GetURL(c, shortenID)
		if err != nil {
			if _, ok := apperror.Is(err, apperror.NotFound); ok {
				c.Redirect(http.StatusSeeOther, handler.defaultURL)
				return
			}

			c.AbortWithStatus(http.StatusNotFound)
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
		log.Println(err)
	}

	c.Redirect(http.StatusSeeOther, url)
}
