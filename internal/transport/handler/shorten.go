package handler

import (
	"cc/internal/domain"
	"cc/internal/dto"
	"cc/internal/service"
	"cc/pkg/base62"
	"cc/pkg/ginutils"
	"cc/pkg/urlutils"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
)

type ShortenHandler struct {
	shortenService service.ShortenService
	authService    service.AuthService
	tagService     service.TagService
	statsService   service.StatsService
}

func NewShortenHandler(
	shortenService service.ShortenService,
	authService service.AuthService,
	tagService service.TagService,
	statsService service.StatsService,
) *ShortenHandler {
	return &ShortenHandler{
		shortenService: shortenService,
		authService:    authService,
		tagService:     tagService,
		statsService:   statsService,
	}
}

func (handler *ShortenHandler) Register(group *gin.RouterGroup) {
	group.GET("/:key/stats", handler.GetShortenStats)
	group.GET("/:key/stats/export", handler.ExportShortenStats)
	group.POST("", handler.CreateShorten)
	group.GET("/:key", handler.GetShorten)
	group.PATCH("/:key", handler.UpdateShorten)
	group.DELETE("/:key", handler.DeleteShorten)
}

func (handler *ShortenHandler) GetShorten(c *gin.Context) {
	shortenKey := c.Param("key")
	shortenID, err := base62.Decode(shortenKey)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var shorten domain.Shorten
	shorten, err = handler.shortenService.GetByID(c, shortenID)
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

	var err error
	request.URL, err = urlutils.Normalize(request.URL)
	if err != nil {
		_ = c.Error(err)
		return
	}

	userID := ginutils.GetUUID(c, "user_id")

	shorten, err := handler.shortenService.Create(c,
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

	userID := ginutils.GetUUID(c, "user_id")

	shortenKey := c.Param("key")
	shortenID, err := base62.Decode(shortenKey)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var shorten domain.Shorten
	shorten, err = handler.shortenService.Update(c,
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

	userID := ginutils.GetUUID(c, "user_id")

	err = handler.shortenService.Delete(c,
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

func (handler *ShortenHandler) ExportShortenStats(c *gin.Context) {
	var request dto.ExportShortenStats
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

	shorten, err := handler.shortenService.GetByID(c, shortenID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	path, err := handler.statsService.ExportStats(c, shorten, request)
	if err != nil {
		_ = c.Error(err)
		return
	}

	_, filename := filepath.Split(path)
	c.FileAttachment(path, filename)

	os.Remove(path)
}
