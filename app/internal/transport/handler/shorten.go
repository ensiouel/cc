package handler

import (
	"cc/app/internal/domain"
	"cc/app/internal/dto"
	"cc/app/internal/service"
	"cc/app/internal/transport"
	"cc/app/internal/transport/middleware"
	"cc/app/pkg/base62"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

type ShortenHandler struct {
	shortenService service.ShortenService
	authService    service.AuthService
	statsService   service.StatsService
	tagService     service.TagService
	cache          *redis.Client
}

func NewShortenHandler(shortenService service.ShortenService, authService service.AuthService, statsService service.StatsService, tagService service.TagService, cache *redis.Client) transport.Handler {
	return &ShortenHandler{shortenService: shortenService, authService: authService, statsService: statsService, tagService: tagService, cache: cache}
}

func (handler *ShortenHandler) Register(group *gin.RouterGroup) {
	group.GET("/:key", handler.Redirect)

	authorized := group.Group("/")
	authorized.Use(middleware.Auth(handler.authService))
	{
		authorized.GET("/api/shortens/:key/stats", handler.GetShortenStats)
		authorized.GET("/api/shortens/:key/stats/export", handler.ExportShortenStats)

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
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	var url string
	url, err = handler.cache.Get(c, "shorten:"+shortenKey).Result()
	if err != nil && errors.Is(err, redis.Nil) == false {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if url == "" {
		url, err = handler.shortenService.GetURL(c, shortenID)
		if err != nil {
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
		c.AbortWithStatus(http.StatusInternalServerError)
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
	shorten, err = handler.shortenService.GetByID(c,
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

	var userID uuid.UUID
	if value, exists := c.Get("user_id"); exists && value != nil {
		userID = value.(uuid.UUID)
	}

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

// ExportShortenStats TODO: move to the service layer
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

	path := filepath.Join(
		"excels",
		fmt.Sprintf("%s_%s_%s",
			shortenKey,
			request.From,
			request.To,
		)+".xlsx",
	)

	var total int64
	total, err = handler.statsService.GetClicksSummary(c,
		shortenID,
		request.From,
		request.To,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	from, _ := time.Parse("2006-01-02", request.From)
	to, _ := time.Parse("2006-01-02", request.To)

	beforeRequest := dto.ExportShortenStats{
		From: from.Add(from.Add(-24 * time.Hour).Sub(to)).Format("2006-01-02"),
		To:   from.Add(-24 * time.Hour).Format("2006-01-02"),
	}

	var totalBefore int64
	totalBefore, err = handler.statsService.GetClicksSummary(c,
		shortenID,
		beforeRequest.From,
		beforeRequest.To,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}
	if err != nil {
		_ = c.Error(err)
		return
	}

	shorten, err := handler.shortenService.GetByID(c, shortenID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	f := excelize.NewFile()
	defer f.Close()

	err = f.SetSheetName("Sheet1", "??????????")
	if err != nil {
		_ = c.Error(err)
		return
	}

	_ = f.SetColWidth("??????????", "A", "E", 32)

	_ = f.SetCellValue("??????????", "A1", "????????????????")
	_ = f.SetCellValue("??????????", "A2", "?????????????????? ????????????")
	_ = f.SetCellValue("??????????", "A3", "???????????????????? ????????????")
	_ = f.SetCellValue("??????????", "A4", "????????????")
	_ = f.SetCellValue("??????????", "A5", "???????????????????????? ????????????")

	_ = f.SetCellValue("??????????", "B1", shorten.Title)
	_ = f.SetCellValue("??????????", "B2", fmt.Sprintf("%s / %s", request.From, request.To))
	_ = f.SetCellValue("??????????", "B3", fmt.Sprintf("%s / %s", beforeRequest.From, beforeRequest.To))
	_ = f.SetCellValue("??????????", "B4", shorten.ShortURL)
	_ = f.SetCellValue("??????????", "B5", shorten.LongURL)

	_ = f.SetCellValue("??????????", "B7", "???????????????????? ????????????")
	_ = f.SetCellValue("??????????", "C7", "?????????????????? ????????????")
	_ = f.SetCellValue("??????????", "D7", "??????????????????")
	_ = f.SetCellValue("??????????", "E7", "?????????????????? %")
	_ = f.SetCellValue("??????????", "A8", "????????????????")

	_ = f.SetCellValue("??????????", "B8", totalBefore)
	_ = f.SetCellValue("??????????", "C8", total)
	_ = f.SetCellFormula("??????????", "D8", "C8-B8")
	_ = f.SetCellFormula("??????????", "E8", "IF(B8 = 0; 100; D8/B8*100)")

	_, _ = f.NewSheet("????????????????")

	_ = f.SetColWidth("????????????????", "A", "D", 32)

	_ = f.SetCellValue("????????????????", "A1", "????????")
	_ = f.SetCellValue("????????????????", "B1", "??????????????????")
	_ = f.SetCellValue("????????????????", "C1", "???????????????????????? ??????????????")
	_ = f.SetCellValue("????????????????", "D1", "???????????????? ????????????????")

	clicks, err := handler.statsService.SelectClicks(c, shortenID, request.From, request.To)
	if err != nil {
		_ = c.Error(err)
		return
	}

	for i, click := range clicks {
		_ = f.SetCellValue("????????????????", fmt.Sprintf("A%d", i+2), click.Timestamp.Format("2006-01-02 15:04"))
		_ = f.SetCellValue("????????????????", fmt.Sprintf("B%d", i+2), click.Platform)
		_ = f.SetCellValue("????????????????", fmt.Sprintf("C%d", i+2), click.OS)
		_ = f.SetCellValue("????????????????", fmt.Sprintf("D%d", i+2), click.Referer)
	}

	err = f.SaveAs(path)
	if err != nil {
		_ = c.Error(err)
		return
	}

	_, filename := filepath.Split(path)
	c.FileAttachment(path, filename)
}
