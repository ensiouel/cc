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
		Referrer:  request.Referer,
		IP:        request.IP,
		Timestamp: request.Timestamp,
	}
	err = service.storage.CreateClick(ctx, clck)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.Internal); ok {
			return apperr.SetScope("create click")
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

func (service *statsService) GetStats(ctx context.Context, shortenID uint64, request dto.GetShortenStats) (stats domain.Stats, err error) {
	var clickMetric model.ClickMetric
	clickMetric, err = service.storage.SelectClickMetric(ctx, shortenID, request.From, request.To, request.Unit, request.Units)
	if err != nil {
		return
	}
	stats.Click = clickMetric.Domain()

	var platformMetrics model.Metrics
	platformMetrics, err = service.storage.SelectPlatformMetrics(ctx, shortenID, request.From, request.To, request.Unit, request.Units)
	if err != nil {
		return
	}
	stats.Platform = platformMetrics.Domain()

	var osMetrics model.Metrics
	osMetrics, err = service.storage.SelectOSMetrics(ctx, shortenID, request.From, request.To, request.Unit, request.Units)
	if err != nil {
		return
	}
	stats.OS = osMetrics.Domain()

	var refererMetrics model.Metrics
	refererMetrics, err = service.storage.SelectRefererMetrics(ctx, shortenID, request.From, request.To, request.Unit, request.Units)
	if err != nil {
		return
	}
	stats.Referer = refererMetrics.Domain()

	return
}
