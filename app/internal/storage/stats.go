package storage

import (
	"cc/app/internal/apperror"
	"cc/app/internal/domain"
	"cc/app/internal/model"
	"cc/app/pkg/postgres"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

const (
	PlatformColumn = "platform"
	OSColumn       = "os"
	RefererColumn  = "referer"
)

type StatsStorage interface {
	CreateClick(ctx context.Context, click model.Click) error

	GetClicksSummary(ctx context.Context, shortenID uint64, from, to string) (int64, error)
	SelectClicks(ctx context.Context, shortenID uint64, from, to string) ([]model.Click, error)

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
    clicks (shorten_id, platform, os, referer, ip, timestamp) 
VALUES 
    ($1, $2, $3, $4, $5, $6)
`

	_, err = storage.client.Exec(ctx, q,
		click.ShortenID,
		click.Platform,
		click.OS,
		click.Referer,
		click.IP,
		click.Timestamp,
	)
	if err != nil {
		return apperror.Internal.WithError(err)
	}

	return
}

func (storage *statsStorage) GetClicksSummary(ctx context.Context, shortenID uint64, from, to string) (total int64, err error) {
	q := `
SELECT
    COUNT(*) as total
FROM clicks
WHERE shorten_id = $1
  AND timestamp BETWEEN $2::TIMESTAMPTZ AND $3::TIMESTAMPTZ + INTERVAL '23 hour 59 minute';
`

	err = storage.client.Get(ctx, &total, q, shortenID, from, to)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return total, apperror.Internal.WithError(err)
	}

	return
}

func (storage *statsStorage) SelectClicks(ctx context.Context, shortenID uint64, from, to string) (clicks []model.Click, err error) {
	q := `
SELECT shorten_id,
       platform,
       os,
       referer,
       ip,
       timestamp
FROM clicks
WHERE shorten_id = $1
  AND timestamp BETWEEN $2::TIMESTAMPTZ AND $3::TIMESTAMPTZ + INTERVAL '23 hour 59 minute'
`

	err = storage.client.Select(ctx, &clicks, q, shortenID, from, to)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return clicks, apperror.Internal.WithError(err)
	}

	return
}

func (storage *statsStorage) SelectClickMetric(ctx context.Context, shortenID uint64, from, to string, unit domain.Unit, units int) (metric model.ClickMetric, err error) {
	q := `
WITH input ("from", "to", unit) AS (VALUES ($2::TIMESTAMPTZ,
                                            $3::TIMESTAMPTZ, $4)),
     series AS (SELECT GENERATE_SERIES(
                               DATE_TRUNC(input.unit, input."from"),
                               DATE_TRUNC(input.unit, input."to" + INTERVAL '23 hour 59 minute'),
                               (1 || input.unit)::INTERVAL
                           ) AS timestamp
                FROM input),
     click AS (SELECT clicks.shorten_id,
                      DATE_TRUNC(input.unit, timestamp) AS timestamp
               FROM clicks,
                    input
               WHERE clicks.shorten_id = $1
                 AND timestamp BETWEEN input."from" AND input."to" + INTERVAL '23 hour 59 minute'),
     metric AS (SELECT COUNT(shorten_id) AS count, series.timestamp AS timestamp
                FROM series
                         LEFT JOIN click ON series.timestamp = click.timestamp
                GROUP BY series.timestamp
                ORDER BY series.timestamp),
     previous AS (SELECT COUNT(*) AS count
                  FROM clicks,
                       input
                  WHERE clicks.shorten_id = $1
                    AND timestamp BETWEEN (input."from" - '1 day'::INTERVAL) - (input."to" - input."from") AND input."from" - INTERVAL '1 day' + INTERVAL '23 hour 59 minute')
SELECT SUM(metric.count)                                                                 AS total,
       SUM(metric.count) - COALESCE(previous.count, 0)                                   AS diff,
       JSON_AGG(JSON_BUILD_OBJECT('timestamp', metric.timestamp, 'count', metric.count)) AS values
FROM metric,
     previous
GROUP BY previous.count
`

	err = storage.client.QueryRow(ctx, q, shortenID, from, to, unit).Scan(
		&metric.Total,
		&metric.Diff,
		&metric.Values,
	)
	if err != nil {
		return metric, apperror.Internal.WithError(err)
	}

	return
}

func (storage *statsStorage) SelectMetrics(ctx context.Context, shortenID uint64, target, from, to string, unit domain.Unit, units int) (metrics []model.Metric, err error) {
	q := `
WITH input ("from", "to", unit) AS (VALUES ($2::TIMESTAMPTZ,
                                                        $3::TIMESTAMPTZ, $4)),
     metric AS (SELECT ` + target + `                          AS name,
                       COUNT(*)                          AS count,
                       DATE_TRUNC(input.unit, timestamp) AS trunc_timestamp
                FROM clicks,
                     input
                WHERE timestamp BETWEEN input."from" AND input."to" + INTERVAL '23 hour 59 minute'
                  AND clicks.shorten_id = $1
                GROUP BY name, trunc_timestamp
                ORDER BY trunc_timestamp),
     previous AS (SELECT ` + target + ` AS name, COUNT(*) AS count
                  FROM clicks,
                       input
                  WHERE clicks.shorten_id = $1
                    AND timestamp BETWEEN (input."from" - INTERVAL '1 day') - (input."to" - input."from") AND input."from" - INTERVAL '1 day' + INTERVAL '23 hour 59 minute'
                  GROUP BY name)
SELECT metric.name                                                                               AS name,
       SUM(metric.count)                                                                         AS total,
       SUM(metric.count) - COALESCE(previous.count, 0)                                           AS diff,
       JSONB_AGG(JSONB_BUILD_OBJECT('timestamp', metric.trunc_timestamp, 'count', metric.count)) AS values
FROM metric
         LEFT JOIN previous ON metric.name = previous.name
GROUP BY metric.name, previous.count;
`

	var rows pgx.Rows
	rows, err = storage.client.Query(ctx, q, shortenID, from, to, unit)
	if err != nil {
		return metrics, apperror.Internal.WithError(err)
	}

	for rows.Next() {
		var metric model.Metric

		err = rows.Scan(&metric.Name, &metric.Total, &metric.Diff, &metric.Values)
		if err != nil {
			return metrics, apperror.Internal.WithError(err)
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
