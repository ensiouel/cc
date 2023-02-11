package service

import (
	"cc/app/internal/apperror"
	"cc/app/internal/domain"
	"cc/app/internal/dto"
	"cc/app/internal/model"
	"cc/app/internal/storage"
	"context"
	"github.com/mileusna/useragent"
	"net/url"
	"time"
)

type StatsService interface {
	CreateClick(ctx context.Context, request dto.CreateClick) error
	CreateClickByUserAgent(ctx context.Context, timestamp time.Time, shortenID uint64, userAgent, referer, ip string) error

	GetClicksSummary(ctx context.Context, shortenID uint64, from, to string) (total int64, err error)
	SelectClicks(ctx context.Context, shortenID uint64, from, to string) ([]domain.Click, error)
	GetStats(ctx context.Context, shortenID uint64, request dto.GetShortenStats) (domain.Stats, error)
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
		if apperr, ok := apperror.Is(err, apperror.TypeInternal); ok {
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

	refererURL, _ := url.Parse(referer)
	if refererURL.Host == "" {
		refererURL.Host = "Other"
	}

	if userAgent.Bot || (platform == "Other" && os == "Other") {
		return nil
	}

	err = service.CreateClick(ctx, dto.CreateClick{
		ShortenID: shortenID,
		Platform:  platform,
		OS:        os,
		Referer:   refererURL.Host,
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
		if apperr, ok := apperror.Is(err, apperror.TypeInternal); ok {
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
		if apperr, ok := apperror.Is(err, apperror.TypeInternal); ok {
			return stats, apperr.WithScope("GetStats.SelectClickMetric")
		}

		return
	}
	stats.Click = clickMetric.Domain()

	var platformMetrics model.Metrics
	platformMetrics, err = service.storage.SelectPlatformMetrics(ctx, shortenID, request.From, request.To, request.Unit, request.Units)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.TypeInternal); ok {
			return stats, apperr.WithScope("GetStats.SelectPlatformMetrics")
		}

		return
	}
	stats.Platform = platformMetrics.Domain()

	var osMetrics model.Metrics
	osMetrics, err = service.storage.SelectOSMetrics(ctx, shortenID, request.From, request.To, request.Unit, request.Units)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.TypeInternal); ok {
			return stats, apperr.WithScope("GetStats.SelectOSMetrics")
		}

		return
	}
	stats.OS = osMetrics.Domain()

	var refererMetrics model.Metrics
	refererMetrics, err = service.storage.SelectRefererMetrics(ctx, shortenID, request.From, request.To, request.Unit, request.Units)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.TypeInternal); ok {
			return stats, apperr.WithScope("GetStats.SelectRefererMetrics")
		}

		return
	}
	stats.Referer = refererMetrics.Domain()

	return
}
