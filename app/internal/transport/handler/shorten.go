package handler

import (
	"cc/app/internal/domain"
	"cc/app/internal/dto"
	"cc/app/internal/service"
	"cc/app/internal/transport"
	"cc/app/internal/transport/middleware/auth"
	"cc/app/pkg/base62"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mileusna/useragent"
	"log"
	"net/http"
	"net/url"
	"time"
)

type ShortenHandler struct {
	shortenService service.ShortenService
	authService    service.AuthService
	statsService   service.StatsService
}

func NewShortenHandler(shortenService service.ShortenService, authService service.AuthService, statsService service.StatsService) transport.Handler {
	return &ShortenHandler{shortenService: shortenService, authService: authService, statsService: statsService}
}

func (handler *ShortenHandler) Register(group *gin.RouterGroup) {
	group.GET("/:key", handler.Redirect)

	authorized := group.Group("/")
	authorized.Use(auth.Middleware(handler.authService))
	{
		authorized.GET("/api/shortens/:key", handler.GetShorten)
		authorized.GET("/api/shortens/:key/:target", handler.GetShortenStats)
		authorized.GET("/api/shortens/:key/:target/summary", handler.GetShortenSummaryStats)
		authorized.POST("/api/shortens", handler.CreateShorten)
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

	var shorten domain.Shorten
	shorten, err = handler.shortenService.GetShortenByID(c, shortenID)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	userAgent := useragent.Parse(c.Request.Header.Get("User-Agent"))

	var platform string

	switch {
	case userAgent.Mobile:
		platform = "Mobile"
	case userAgent.Desktop:
		platform = "Desktop"
	case userAgent.Tablet:
		platform = "Tablet"
	default:
		platform = "Unknown"
	}

	var os string
	os = userAgent.OS
	if os == "" {
		os = "Unknown"
	}

	referrer, _ := url.Parse(c.Request.Referer())
	if referrer.Host == "" {
		referrer.Host = "Unknown"
	}

	if !userAgent.Bot {
		err = handler.statsService.CreateClick(c, dto.CreateClick{
			ShortenID: shortenID,
			Platform:  platform,
			OS:        os,
			Referrer:  referrer.Host,
			Timestamp: time.Now(),
		})
		if err != nil {
			_ = c.Error(err)
			return
		}
	}

	c.Redirect(http.StatusSeeOther, shorten.LongURL)
}

func (handler *ShortenHandler) GetShorten(c *gin.Context) {
	shortenID, err := base62.Decode(c.Param("key"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	var shorten domain.Shorten
	shorten, err = handler.shortenService.GetShortenByID(c, shortenID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"response": shorten,
	})
}

func (handler *ShortenHandler) GetShortenStats(c *gin.Context) {
	var request dto.GetShortenStats

	if err := c.BindQuery(&request); err != nil {
		log.Println(2 + 2)
		_ = c.Error(err)
		return
	}

	if err := request.Validate(); err != nil {
		_ = c.Error(err)
		return
	}

	shortenID, err := base62.Decode(c.Param("key"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	target := c.Param("target")

	var stats any

	switch target {
	case "click":
		stats, err = handler.statsService.GetClickStats(c, shortenID, request)
		if err != nil {
			_ = c.Error(err)
			return
		}
	case "platform", "referrer", "os":
		stats, err = handler.statsService.GetMetricStats(c, target, shortenID, request)
		if err != nil {
			_ = c.Error(err)
			return
		}
	default:
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"response": stats,
	})
}

func (handler *ShortenHandler) GetShortenSummaryStats(c *gin.Context) {
	var request dto.GetShortenSummaryStats

	if err := c.BindQuery(&request); err != nil {
		_ = c.Error(err)
		return
	}

	if err := request.Validate(); err != nil {
		_ = c.Error(err)
		return
	}

	shortenID, err := base62.Decode(c.Param("key"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	target := c.Param("target")

	var stats any

	switch target {
	case "click":
		stats, err = handler.statsService.GetClickSummaryStats(c, shortenID, request)
		if err != nil {
			_ = c.Error(err)
			return
		}
	case "platform", "referrer", "os":
		stats, err = handler.statsService.GetMetricSummaryStats(c, target, shortenID, request)
		if err != nil {
			_ = c.Error(err)
			return
		}
	default:
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"response": stats,
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
	if value, exists := c.Get("user_id"); exists {
		userID = value.(uuid.UUID)
	}

	shorten, err := handler.shortenService.CreateShorten(c, userID, request)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
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
	if value, exists := c.Get("user_id"); exists {
		userID = value.(uuid.UUID)
	}

	shortenID, err := base62.Decode(c.Param("key"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	var shorten domain.Shorten
	shorten, err = handler.shortenService.UpdateShorten(c, userID, shortenID, request)
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
	if value, exists := c.Get("user_id"); exists {
		userID = value.(uuid.UUID)
	}

	err = handler.shortenService.DeleteShorten(c, userID, shortenID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": 1,
	})
}
