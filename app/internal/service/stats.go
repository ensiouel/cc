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

	GetClickStats(ctx context.Context, shortenID uint64, request dto.GetShortenStats) (domain.ClickStats, error)
	GetUniqueClickStats(ctx context.Context, shortenID uint64, request dto.GetShortenStats) (domain.ClickStats, error)

	GetMetricStats(ctx context.Context, target string, shortenID uint64, request dto.GetShortenStats) (domain.MetricStats, error)
	GetSummaryMetricStats(ctx context.Context, target string, shortenID uint64, request dto.GetShortenSummaryStats) (domain.SummaryMetricStats, error)
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
		IP:        request.IP,
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
		platform = "Unknown"
	}

	os = userAgent.OS
	if os == "" {
		os = "Unknown"
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
		IP:        ip,
		Timestamp: timestamp,
	})
	if err != nil {
		return err
	}

	return nil
}

func (service *statsService) GetClickStats(ctx context.Context, shortenID uint64, request dto.GetShortenStats) (stats domain.ClickStats, err error) {
	var (
		from, to time.Time
	)
	from, err = time.Parse("2006-01-02", request.From)
	to, err = time.Parse("2006-01-02", request.To)
	if err != nil {
		return stats, apperror.ErrInvalidParams
	}

	var clickStats model.ClickStats
	clickStats, err = service.storage.SelectClicks(ctx, shortenID, from, to, string(request.Unit), request.Units)
	if err != nil {
		if apperr, ok := apperror.Internal(err); ok {
			return stats, apperr.SetScope("get click stats")
		}

		return
	}

	var clicks int
	clicks, err = service.storage.GetTotalClicks(ctx, shortenID, from, to)
	if err != nil {
		if apperr, ok := apperror.Internal(err); ok {
			return stats, apperr.SetScope("get click stats")
		}

		return
	}

	stats = domain.ClickStats{
		Total:  clicks,
		Clicks: clickStats.Domain(),
		Unit:   request.Unit,
		Units:  len(clickStats),
	}

	return
}

func (service *statsService) GetUniqueClickStats(ctx context.Context, shortenID uint64, request dto.GetShortenStats) (stats domain.ClickStats, err error) {
	var (
		from, to time.Time
	)
	from, err = time.Parse("2006-01-02", request.From)
	to, err = time.Parse("2006-01-02", request.To)
	if err != nil {
		return stats, apperror.ErrInvalidParams
	}

	var clickStats model.ClickStats
	clickStats, err = service.storage.SelectUniqueClicks(ctx, shortenID, from, to, string(request.Unit), request.Units)
	if err != nil {
		if apperr, ok := apperror.Internal(err); ok {
			return stats, apperr.SetScope("get click stats")
		}

		return
	}

	var clicks int
	clicks, err = service.storage.GetTotalUniqueClicks(ctx, shortenID, from, to, string(request.Unit))
	if err != nil {
		if apperr, ok := apperror.Internal(err); ok {
			return stats, apperr.SetScope("get click stats")
		}

		return
	}

	stats = domain.ClickStats{
		Total:  clicks,
		Clicks: clickStats.Domain(),
		Unit:   request.Unit,
		Units:  len(clickStats),
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

func (service *statsService) GetSummaryMetricStats(ctx context.Context, target string, shortenID uint64, request dto.GetShortenSummaryStats) (metricSummaryStats domain.SummaryMetricStats, err error) {
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

	metricSummaryStats = domain.SummaryMetricStats{
		Metrics: metrics.Domain(),
	}

	return
}
