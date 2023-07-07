package service

import (
	"cc/internal/domain"
	"cc/internal/dto"
	"cc/internal/model"
	"cc/internal/storage"
	"cc/pkg/apperror"
	"cc/pkg/base62"
	"context"
	"fmt"
	"github.com/goware/urlx"
	"github.com/mileusna/useragent"
	"github.com/xuri/excelize/v2"
	"path/filepath"
	"strings"
	"time"
)

type StatsService interface {
	CreateClick(ctx context.Context, request dto.CreateClick) error
	CreateClickByUserAgent(ctx context.Context, timestamp time.Time, shortenID uint64, userAgent, referer, ip string) error
	GetClicksSummary(ctx context.Context, shortenID uint64, from, to string) (total int64, err error)
	SelectClicks(ctx context.Context, shortenID uint64, from, to string) ([]domain.Click, error)
	GetStats(ctx context.Context, shortenID uint64, request dto.GetShortenStats) (domain.Stats, error)
	ExportStats(ctx context.Context, shorten domain.Shorten, request dto.ExportShortenStats) (string, error)
}

type statsService struct {
	storage storage.StatsStorage
}

func NewStatsService(storage storage.StatsStorage) StatsService {
	return &statsService{storage: storage}
}

func (service *statsService) CreateClick(ctx context.Context, request dto.CreateClick) (err error) {
	clck := model.Click{
		ShortenID: request.ShortenID,
		Platform:  request.Platform,
		OS:        request.OS,
		Referer:   request.Referer,
		IP:        request.IP,
		Timestamp: request.Timestamp,
	}
	err = service.storage.CreateClick(ctx, clck)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.Internal); ok {
			return apperr.WithScope("create click")
		}

		return
	}

	return
}

func (service *statsService) CreateClickByUserAgent(ctx context.Context, timestamp time.Time, shortenID uint64, ua, referer, ip string) (err error) {
	userAgent := useragent.Parse(ua)

	var platform, os string
	switch {
	case userAgent.Mobile:
		platform = "Mobile"
	case userAgent.Desktop:
		platform = "Desktop"
	case userAgent.Tablet:
		platform = "Tablet"
	default:
		platform = "Other"
	}

	os = userAgent.OS
	if os == "" {
		os = "Other"
	}

	if referer == "" {
		referer = "Other"
	} else {
		referer, _ = urlx.NormalizeString(referer)
		referer = strings.Replace(referer, "www.", "", 1)

		parse, err := urlx.Parse(referer)
		if err != nil {
			return
		}
		parse.RawQuery = ""

		referer = parse.String()
	}

	if userAgent.Bot || (platform == "Other" && os == "Other") {
		return nil
	}

	err = service.CreateClick(ctx, dto.CreateClick{
		ShortenID: shortenID,
		Platform:  platform,
		OS:        os,
		Referer:   referer,
		IP:        ip,
		Timestamp: timestamp,
	})
	if err != nil {
		return err
	}

	return nil
}

func (service *statsService) GetClicksSummary(ctx context.Context, shortenID uint64, from, to string) (total int64, err error) {
	return service.storage.GetClicksSummary(ctx, shortenID, from, to)
}

func (service *statsService) SelectClicks(ctx context.Context, shortenID uint64, from, to string) (clicks []domain.Click, err error) {
	var clcks model.Clicks
	clcks, err = service.storage.SelectClicks(ctx, shortenID, from, to)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.Internal); ok {
			return clicks, apperr.WithScope("SelectClicks")
		}

		return
	}

	return clcks.Domain(), nil
}

func (service *statsService) GetStats(ctx context.Context, shortenID uint64, request dto.GetShortenStats) (stats domain.Stats, err error) {
	var clickMetric model.ClickMetric
	clickMetric, err = service.storage.SelectClickMetric(ctx, shortenID, request.From, request.To, request.Unit, request.Units)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.Internal); ok {
			return stats, apperr.WithScope("GetStats.SelectClickMetric")
		}

		return
	}
	stats.Click = clickMetric.Domain()

	var platformMetrics model.Metrics
	platformMetrics, err = service.storage.SelectPlatformMetrics(ctx, shortenID, request.From, request.To, request.Unit, request.Units)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.Internal); ok {
			return stats, apperr.WithScope("GetStats.SelectPlatformMetrics")
		}

		return
	}
	stats.Platform = platformMetrics.Domain()

	var osMetrics model.Metrics
	osMetrics, err = service.storage.SelectOSMetrics(ctx, shortenID, request.From, request.To, request.Unit, request.Units)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.Internal); ok {
			return stats, apperr.WithScope("GetStats.SelectOSMetrics")
		}

		return
	}
	stats.OS = osMetrics.Domain()

	var refererMetrics model.Metrics
	refererMetrics, err = service.storage.SelectRefererMetrics(ctx, shortenID, request.From, request.To, request.Unit, request.Units)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.Internal); ok {
			return stats, apperr.WithScope("GetStats.SelectRefererMetrics")
		}

		return
	}
	stats.Referer = refererMetrics.Domain()

	return
}

func (service *statsService) ExportStats(ctx context.Context, shorten domain.Shorten, request dto.ExportShortenStats) (string, error) {
	shortenID, err := base62.Decode(shorten.ID)
	if err != nil {
		return "", err
	}

	total, err := service.GetClicksSummary(ctx,
		shortenID,
		request.From,
		request.To,
	)
	if err != nil {
		return "", err
	}

	from, _ := time.Parse("2006-01-02", request.From)
	to, _ := time.Parse("2006-01-02", request.To)

	beforeRequest := dto.ExportShortenStats{
		From: from.Add(from.Add(-24 * time.Hour).Sub(to)).Format("2006-01-02"),
		To:   from.Add(-24 * time.Hour).Format("2006-01-02"),
	}

	var totalBefore int64
	totalBefore, err = service.GetClicksSummary(ctx,
		shortenID,
		beforeRequest.From,
		beforeRequest.To,
	)
	if err != nil {
		return "", err
	}

	f := excelize.NewFile()
	defer f.Close()

	err = f.SetSheetName("Sheet1", "Обзор")
	if err != nil {
		return "", err
	}

	_ = f.SetColWidth("Обзор", "A", "E", 32)

	_ = f.SetCellValue("Обзор", "A1", "Название")
	_ = f.SetCellValue("Обзор", "A2", "Выбранный период")
	_ = f.SetCellValue("Обзор", "A3", "Предыдущий период")
	_ = f.SetCellValue("Обзор", "A4", "Ссылка")
	_ = f.SetCellValue("Обзор", "A5", "Оригинальная ссылка")

	_ = f.SetCellValue("Обзор", "B1", shorten.Title)
	_ = f.SetCellValue("Обзор", "B2", fmt.Sprintf("%s / %s", request.From, request.To))
	_ = f.SetCellValue("Обзор", "B3", fmt.Sprintf("%s / %s", beforeRequest.From, beforeRequest.To))
	_ = f.SetCellValue("Обзор", "B4", shorten.ShortURL)
	_ = f.SetCellValue("Обзор", "B5", shorten.LongURL)

	_ = f.SetCellValue("Обзор", "B7", "Предыдущий период")
	_ = f.SetCellValue("Обзор", "C7", "Выбранный период")
	_ = f.SetCellValue("Обзор", "D7", "Изменение")
	_ = f.SetCellValue("Обзор", "E7", "Изменение %")
	_ = f.SetCellValue("Обзор", "A8", "Переходы")

	_ = f.SetCellValue("Обзор", "B8", totalBefore)
	_ = f.SetCellValue("Обзор", "C8", total)
	_ = f.SetCellFormula("Обзор", "D8", "C8-B8")
	_ = f.SetCellFormula("Обзор", "E8", "IF(B8 = 0; 100; D8/B8*100)")

	_, _ = f.NewSheet("Переходы")

	_ = f.SetColWidth("Переходы", "A", "D", 32)

	_ = f.SetCellValue("Переходы", "A1", "Дата")
	_ = f.SetCellValue("Переходы", "B1", "Платформа")
	_ = f.SetCellValue("Переходы", "C1", "Операционная система")
	_ = f.SetCellValue("Переходы", "D1", "Источник перехода")

	clicks, err := service.SelectClicks(ctx, shortenID, request.From, request.To)
	if err != nil {
		return "", err
	}

	for i, click := range clicks {
		_ = f.SetCellValue("Переходы", fmt.Sprintf("A%d", i+2), click.Timestamp.Format("2006-01-02 15:04"))
		_ = f.SetCellValue("Переходы", fmt.Sprintf("B%d", i+2), click.Platform)
		_ = f.SetCellValue("Переходы", fmt.Sprintf("C%d", i+2), click.OS)
		_ = f.SetCellValue("Переходы", fmt.Sprintf("D%d", i+2), click.Referer)
	}

	path := filepath.Join(
		fmt.Sprintf("%s_%s_%s_%d",
			shorten.ID,
			from.Format("20060102"),
			to.Format("20060102"),
			time.Now().Round(1*time.Hour).Unix(),
		) + ".xlsx",
	)

	err = f.SaveAs(path)
	if err != nil {
		return "", err
	}

	return path, nil
}
