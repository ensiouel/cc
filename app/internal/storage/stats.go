package storage

import (
	"cc/app/internal/apperror"
	"cc/app/internal/model"
	"cc/app/pkg/postgres"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"time"
)

/*

SelectClicks
{
	"total": 1000,
	"clicks": [
		{
			"count": 2,
			"date": ""
		},
		{
			"count": 4,
			"date": ""
		},
	]
}

SelectUniqueClicks
{
	"total": 998,
	"clicks": [
		{
			"count": 2,
			"date": ""
		},
		{
			"count": 2,
			"date": ""
		},
	]
}

*/

type StatsStorage interface {
	CreateClick(ctx context.Context, click model.Click) error

	SelectUniqueClicks(ctx context.Context, shortenID uint64, from, to time.Time, unit string, units int) ([]model.ClickState, error)
	SelectClicks(ctx context.Context, shortenID uint64, from, to time.Time, unit string, units int) ([]model.ClickState, error)

	GetTotalClicks(ctx context.Context, shortenID uint64, from, to time.Time) (int, error)
	GetTotalUniqueClicks(ctx context.Context, shortenID uint64, from, to time.Time, unit string) (int, error)

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

	//TODO fix it
	err = storage.client.Select(ctx, &stats, q, from.Format("2006-01-02"), to.Format("2006-01-02"), unit, shortenID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return stats, apperror.ErrNotExists
		}

		return stats, apperror.ErrInternalError.SetError(err)
	}

	return
}

func (storage *statsStorage) SelectUniqueClicks(ctx context.Context, shortenID uint64, from, to time.Time, unit string, units int) (stats []model.ClickState, err error) {
	q := `
WITH RANGE (f, t, i) as (
    VALUES ($1::TIMESTAMPTZ, $2::TIMESTAMPTZ + INTERVAL '23 hour 59 minute', $3)
),
new_clicks AS (
    SELECT DISTINCT ON (ip, date_trunc(RANGE.i, timestamp))
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
            (SELECT MIN(date_trunc(RANGE.i, timestamp)) FROM new_clicks),
            (SELECT MAX(date_trunc(RANGE.i, timestamp)) FROM new_clicks),
            (1 || RANGE.i)::interval
        )
)
SELECT
    COUNT(shorten_id) as count,
    RANGE_SERIES.timestamp as date
FROM
    RANGE_SERIES
LEFT JOIN
    new_clicks
ON
    RANGE_SERIES.timestamp = new_clicks.timestamp
GROUP BY
    RANGE_SERIES.timestamp
ORDER BY
    RANGE_SERIES.timestamp;
`

	err = storage.client.Select(ctx, &stats, q, from.Format("2006-01-02"), to.Format("2006-01-02"), unit, shortenID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return stats, apperror.ErrNotExists
		}

		return stats, apperror.ErrInternalError.SetError(err)
	}

	return
}

func (storage *statsStorage) GetTotalClicks(ctx context.Context, shortenID uint64, from, to time.Time) (clicks int, err error) {
	q := `
SELECT 
    COUNT(shorten_id)
FROM 
    clicks
WHERE 
    shorten_id = $1 AND
    timestamp BETWEEN $2::TIMESTAMPTZ AND $3::TIMESTAMPTZ + INTERVAL '23 hour 59 minute'
`

	err = storage.client.Get(ctx, &clicks, q, shortenID, from.Format("2006-01-02"), to.Format("2006-01-02"))
	if err != nil {
		return
	}

	return
}

func (storage *statsStorage) GetTotalUniqueClicks(ctx context.Context, shortenID uint64, from, to time.Time, unit string) (clicks int, err error) {
	q := `
SELECT
    COUNT(DISTINCT (ip, date_trunc($4, timestamp)))
FROM
    clicks
WHERE
    shorten_id = $1 AND
    timestamp BETWEEN $2::TIMESTAMPTZ AND $3::TIMESTAMPTZ + INTERVAL '23 hour 59 minute'
`

	err = storage.client.Get(ctx, &clicks, q, shortenID, from.Format("2006-01-02"), to.Format("2006-01-02"), unit)
	if err != nil {
		return
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

	err = storage.client.Select(ctx, &stats, q, from.Format("2006-01-02"), to.Format("2006-01-02"), unit, shortenID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return stats, apperror.ErrNotExists
		}

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

	err = storage.client.Select(ctx, &stats, q, from.Format("2006-01-02"), to.Format("2006-01-02"), shortenID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return stats, apperror.ErrNotExists
		}

		return stats, apperror.ErrInternalError.SetError(err)
	}

	return
}
