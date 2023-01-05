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
	CreateClickByUserAgent(ctx context.Context, shortenID uint64, userAgent, referer string) error
	GetClickStats(ctx context.Context, shortenID uint64, request dto.GetShortenStats) (domain.ClickStats, error)
	GetClickSummaryStats(ctx context.Context, shortenID uint64, request dto.GetShortenSummaryStats) (domain.ClickSummaryStats, error)
	GetMetricStats(ctx context.Context, target string, shortenID uint64, request dto.GetShortenStats) (domain.MetricStats, error)
	GetMetricSummaryStats(ctx context.Context, target string, shortenID uint64, request dto.GetShortenSummaryStats) (domain.MetricSummaryStats, error)
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
		Referrer:  request.Referrer,
		Timestamp: request.Timestamp,
	}
	err = service.storage.CreateClick(ctx, clck)
	if err != nil {
		if apperr, ok := apperror.Internal(err); ok {
			return apperr.SetScope("create click")
		}

		return
	}

	return
}

func (service *statsService) CreateClickByUserAgent(ctx context.Context, shortenID uint64, ua, referer string) (err error) {
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

	referrer, _ := url.Parse(referer)
	if referrer.Host == "" {
		referrer.Host = "Unknown"
	}

	err = service.CreateClick(ctx, dto.CreateClick{
		ShortenID: shortenID,
		Platform:  platform,
		OS:        os,
		Referrer:  referrer.Host,
		Timestamp: time.Now(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (service *statsService) GetClickStats(ctx context.Context, shortenID uint64, request dto.GetShortenStats) (clickStats domain.ClickStats, err error) {
	var (
		from, to time.Time
	)
	from, err = time.Parse("2006-01-02", request.From)
	to, err = time.Parse("2006-01-02", request.To)
	if err != nil {
		return clickStats, apperror.ErrInvalidParams
	}

	var clicks model.ClickStats
	clicks, err = service.storage.SelectClicks(ctx, shortenID, from, to, string(request.Unit), request.Units)
	if err != nil {
		if apperr, ok := apperror.Internal(err); ok {
			return clickStats, apperr.SetScope("get click stats")
		}

		return
	}

	clickStats = domain.ClickStats{
		Clicks: clicks.Domain(),
		Unit:   request.Unit,
		Units:  len(clicks),
	}

	return
}

func (service *statsService) GetClickSummaryStats(ctx context.Context, shortenID uint64, request dto.GetShortenSummaryStats) (clickSummaryStats domain.ClickSummaryStats, err error) {
	var (
		from, to time.Time
	)
	from, err = time.Parse("2006-01-02", request.From)
	to, err = time.Parse("2006-01-02", request.To)
	if err != nil {
		return clickSummaryStats, apperror.ErrInvalidParams
	}

	var clickSummary int
	clickSummary, err = service.storage.SelectSummaryClicks(ctx, shortenID, from, to)
	if err != nil {
		if apperr, ok := apperror.Internal(err); ok {
			return clickSummaryStats, apperr.SetScope("get click summary stats")
		}

		return
	}

	clickSummaryStats = domain.ClickSummaryStats{
		Clicks: clickSummary,
	}

	return
}

func (service *statsService) GetMetricStats(ctx context.Context, target string, shortenID uint64, request dto.GetShortenStats) (metricStats domain.MetricStats, err error) {
	var (
		from, to time.Time
	)
	from, err = time.Parse("2006-01-02", request.From)
	to, err = time.Parse("2006-01-02", request.To)
	if err != nil {
		return metricStats, apperror.ErrInvalidParams
	}

	var metrics model.MetricStats
	metrics, err = service.storage.SelectMetrics(ctx, target, shortenID, from, to, string(request.Unit), request.Units)
	if err != nil {
		if apperr, ok := apperror.Internal(err); ok {
			return metricStats, apperr.SetScope("get metric stats")
		}

		return
	}

	metricStats = domain.MetricStats{
		Metrics: metrics.Domain(),
		Unit:    request.Unit,
		Units:   len(metrics),
	}

	return
}

func (service *statsService) GetMetricSummaryStats(ctx context.Context, target string, shortenID uint64, request dto.GetShortenSummaryStats) (metricSummaryStats domain.MetricSummaryStats, err error) {
	var (
		from, to time.Time
	)
	from, err = time.Parse("2006-01-02", request.From)
	to, err = time.Parse("2006-01-02", request.To)
	if err != nil {
		return metricSummaryStats, apperror.ErrInvalidParams
	}

	var metrics model.MetricSummaryStats
	metrics, err = service.storage.SelectSummaryMetrics(ctx, target, shortenID, from, to)
	if err != nil {
		if apperr, ok := apperror.Internal(err); ok {
			return metricSummaryStats, apperr.SetScope("get metric summary stats")
		}

		return
	}

	metricSummaryStats = domain.MetricSummaryStats{
		Metrics: metrics.Domain(),
	}

	return
}
