package storage

import (
	"cc/app/internal/apperror"
	"cc/app/internal/model"
	"cc/app/pkg/postgres"
	"context"
	"time"
)

type StatsStorage interface {
	CreateClick(ctx context.Context, click model.Click) error
	SelectClicks(ctx context.Context, shortenID uint64, from, to time.Time, unit string, units int) ([]model.ClickState, error)
	SelectSummaryClicks(ctx context.Context, shortenID uint64, from, to time.Time) (int, error)
	SelectMetrics(ctx context.Context, target string, shortenID uint64, from, to time.Time, unit string, units int) ([]model.MetricState, error)
	SelectSummaryMetrics(ctx context.Context, target string, shortenID uint64, from, to time.Time) ([]model.MetricSummaryState, error)
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
    clicks (shorten_id, platform, os, referrer, timestamp) 
VALUES 
    ($1, $2, $3, $4, $5)
`

	_, err = storage.client.Exec(ctx, q,
		click.ShortenID,
		click.Platform,
		click.OS,
		click.Referrer,
		click.Timestamp,
	)
	if err != nil {
		return apperror.ErrInternalError.SetError(err)
	}

	return
}

func (storage *statsStorage) SelectClicks(ctx context.Context, shortenID uint64, from, to time.Time, unit string, units int) (stats []model.ClickState, err error) {
	q := `
WITH RANGE (f, t, i) as (
    VALUES ($1::TIMESTAMPTZ, $2::TIMESTAMPTZ + INTERVAL '23 hour 59 minute', $3)
),
clicks AS (
    SELECT
    	shorten_id, date_trunc(RANGE.i, timestamp) as timestamp
    FROM
        clicks, RANGE
    WHERE
        shorten_id = $4 AND
        timestamp BETWEEN RANGE.f AND RANGE.t
),
RANGE_SERIES AS (
    SELECT
        GENERATE_SERIES AS timestamp
    FROM
        RANGE,
        GENERATE_SERIES (
            (SELECT MIN(date_trunc(RANGE.i, timestamp)) FROM clicks),
            (SELECT MAX(date_trunc(RANGE.i, timestamp)) FROM clicks),
            (1 || RANGE.i)::interval
        )
)
SELECT
    RANGE_SERIES.timestamp as date,
    COUNT(shorten_id) as count
FROM
    RANGE_SERIES
LEFT JOIN
    clicks
ON
	RANGE_SERIES.timestamp = clicks.timestamp
GROUP BY
    RANGE_SERIES.timestamp
ORDER BY
    RANGE_SERIES.timestamp;
`

	err = storage.client.Select(ctx, &stats, q, from, to, unit, shortenID)
	if err != nil {
		return stats, apperror.ErrInternalError.SetError(err)
	}

	return
}

func (storage *statsStorage) SelectSummaryClicks(ctx context.Context, shortenID uint64, from, to time.Time) (stats int, err error) {
	q := `
WITH RANGE (f, t) as (
    values ($1::TIMESTAMPTZ, $2::TIMESTAMPTZ + INTERVAL '23 hour 59 minute')
)
SELECT
    COUNT(shorten_id) as count
FROM
    clicks, RANGE
WHERE
    shorten_id = $3 AND
    timestamp BETWEEN RANGE.f AND RANGE.t
GROUP BY shorten_id;
`

	err = storage.client.Get(ctx, &stats, q, from, to, shortenID)
	if err != nil {
		return stats, apperror.ErrInternalError.SetError(err)
	}

	return
}

func (storage *statsStorage) SelectMetrics(ctx context.Context, target string, shortenID uint64, from, to time.Time, unit string, units int) (stats []model.MetricState, err error) {
	q := `
WITH RANGE (f, t, i) as (
    VALUES ($1::TIMESTAMPTZ, $2::TIMESTAMPTZ + INTERVAL '23 hour 59 minute', $3)
),
select_clicks AS (
    SELECT
        shorten_id, ` + target + ` as name, date_trunc(RANGE.i, timestamp) as timestamp
    FROM
        clicks, RANGE
    WHERE
        shorten_id = $4 AND
        timestamp BETWEEN RANGE.f AND RANGE.t
)
SELECT
    name,
    COUNT(shorten_id) as count,
    timestamp as date
FROM
    select_clicks
GROUP BY
    shorten_id, timestamp, name
ORDER BY
    timestamp;
`

	err = storage.client.Select(ctx, &stats, q, from, to, unit, shortenID)
	if err != nil {
		return stats, apperror.ErrInternalError.SetError(err)
	}

	return
}

func (storage *statsStorage) SelectSummaryMetrics(ctx context.Context, target string, shortenID uint64, from, to time.Time) (stats []model.MetricSummaryState, err error) {
	q := `
WITH RANGE (f, t) as (
    values ($1::TIMESTAMPTZ, $2::TIMESTAMPTZ + INTERVAL '23 hour 59 minute')
)
SELECT
    ` + target + ` as name, COUNT(shorten_id) as count
FROM
    clicks, RANGE
WHERE
    shorten_id = $3 AND
    timestamp BETWEEN RANGE.f AND RANGE.t
GROUP BY shorten_id, name;
`

	err = storage.client.Select(ctx, &stats, q, from, to, shortenID)
	if err != nil {
		return stats, apperror.ErrInternalError.SetError(err)
	}

	return
}
