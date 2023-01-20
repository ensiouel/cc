package storage

import (
	"cc/app/internal/apperror"
	"cc/app/internal/domain"
	"cc/app/internal/model"
	"cc/app/pkg/postgres"
	"context"
	"github.com/jackc/pgx/v5"
)

const (
	PlatformColumn = "platform"
	OSColumn       = "os"
	RefererColumn  = "referer"
)

type StatsStorage interface {
	CreateClick(ctx context.Context, click model.Click) error

	SelectClickMetric(ctx context.Context, shortenID uint64, from, to string, unit domain.Unit, units int) (model.ClickMetric, error)
	SelectMetrics(ctx context.Context, shortenID uint64, target, from, to string, unit domain.Unit, units int) ([]model.Metric, error)

	SelectPlatformMetrics(ctx context.Context, shortenID uint64, from, to string, unit domain.Unit, units int) ([]model.Metric, error)
	SelectOSMetrics(ctx context.Context, shortenID uint64, from, to string, unit domain.Unit, units int) ([]model.Metric, error)
	SelectRefererMetrics(ctx context.Context, shortenID uint64, from, to string, unit domain.Unit, units int) ([]model.Metric, error)
}

type statsStorage struct {
	client postgres.Client
}

func NewStatsStorage(client postgres.Client) StatsStorage {
	return &statsStorage{client: client}
}

func (storage *statsStorage) CreateClick(ctx context.Context, click model.Click) (err error) {
	q := `
INSERT INTO 
    clicks (shorten_id, platform, os, referrer, ip, timestamp) 
VALUES 
    ($1, $2, $3, $4, $5, $6)
`

	_, err = storage.client.Exec(ctx, q,
		click.ShortenID,
		click.Platform,
		click.OS,
		click.Referrer,
		click.IP,
		click.Timestamp,
	)
	if err != nil {
		return apperror.ErrInternalError.SetError(err)
	}

	return
}

func (storage *statsStorage) SelectClickMetric(ctx context.Context, shortenID uint64, from, to string, unit domain.Unit, units int) (metric model.ClickMetric, err error) {
	q := `
WITH series AS (SELECT GENERATE_SERIES(
                               DATE_TRUNC($3, $1::TIMESTAMPTZ),
                               DATE_TRUNC($3, $2::TIMESTAMPTZ + INTERVAL '23 hour 59 minute'),
                               (1 || $3)::INTERVAL
                           ) AS timestamp),
     click AS (SELECT shorten_id,
                      DATE_TRUNC($3, timestamp) as timestamp
               FROM clicks
               WHERE shorten_id = $4
                 AND timestamp BETWEEN $1::TIMESTAMPTZ AND $2::TIMESTAMPTZ + INTERVAL '23 hour 59 minute'),
     metric AS (SELECT COUNT(shorten_id) AS count, series.timestamp AS timestamp
                FROM series
                         LEFT JOIN click ON series.timestamp = click.timestamp
                GROUP BY series.timestamp
                ORDER BY series.timestamp)
SELECT SUM(count)                                                                        as total,
       JSON_AGG(JSON_BUILD_OBJECT('timestamp', metric.timestamp, 'count', metric.count)) as values
FROM metric;`

	err = storage.client.QueryRow(ctx, q, from, to, unit, shortenID).Scan(
		&metric.Total,
		&metric.Values,
	)
	if err != nil {
		return metric, apperror.ErrInternalError.SetError(err)
	}

	return
}

func (storage *statsStorage) SelectMetrics(ctx context.Context, shortenID uint64, target, from, to string, unit domain.Unit, units int) (metrics []model.Metric, err error) {
	q := `
WITH metric AS (SELECT ` + target + `                     as name,
                       COUNT(*)                     as count,
                       DATE_TRUNC($1, timestamp) as trunc_timestamp
                FROM clicks
                WHERE timestamp BETWEEN $2::TIMESTAMPTZ AND $3::TIMESTAMPTZ + INTERVAL '23 hour 59 minute'
                  AND shorten_id = $4
                GROUP BY name, trunc_timestamp
                ORDER BY trunc_timestamp)
SELECT metric.name                                                                             as name,
       SUM(metric.count)                                                                       as total,
       JSON_AGG(JSON_BUILD_OBJECT('timestamp', metric.trunc_timestamp, 'count', metric.count)) as values
FROM metric
GROUP BY metric.name;
`

	var rows pgx.Rows
	rows, err = storage.client.Query(ctx, q, unit, from, to, shortenID)
	if err != nil {
		return metrics, apperror.ErrInternalError.SetError(err)
	}

	for rows.Next() {
		var metric model.Metric

		err = rows.Scan(&metric.Name, &metric.Total, &metric.Values)
		if err != nil {
			return metrics, apperror.ErrInternalError.SetError(err)
		}

		metrics = append(metrics, metric)
	}

	return
}

func (storage *statsStorage) SelectPlatformMetrics(ctx context.Context, shortenID uint64, from, to string, unit domain.Unit, units int) ([]model.Metric, error) {
	return storage.SelectMetrics(ctx, shortenID, PlatformColumn, from, to, unit, units)
}

func (storage *statsStorage) SelectOSMetrics(ctx context.Context, shortenID uint64, from, to string, unit domain.Unit, units int) ([]model.Metric, error) {
	return storage.SelectMetrics(ctx, shortenID, OSColumn, from, to, unit, units)
}

func (storage *statsStorage) SelectRefererMetrics(ctx context.Context, shortenID uint64, from, to string, unit domain.Unit, units int) ([]model.Metric, error) {
	return storage.SelectMetrics(ctx, shortenID, RefererColumn, from, to, unit, units)
}
