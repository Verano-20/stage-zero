package telemetry

import (
	"context"
	"database/sql"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type AppMetrics struct {
	// HTTP metrics
	HTTPRequestsTotal   metric.Int64Counter
	HTTPRequestDuration metric.Float64Histogram
	HTTPRequestsActive  metric.Int64UpDownCounter

	// Database metrics
	DBConnectionsOpen  metric.Int64Gauge
	DBConnectionsInUse metric.Int64Gauge
	DBConnectionsIdle  metric.Int64Gauge
	DBQueriesTotal     metric.Int64Counter
	DBQueryDuration    metric.Float64Histogram

	// Auth metrics
	AuthAttemptsTotal metric.Int64Counter
	AuthFailuresTotal metric.Int64Counter

	// Business metrics
	UsersTotal   metric.Int64UpDownCounter
	SimplesTotal metric.Int64UpDownCounter
}

func NewAppMetrics(meter metric.Meter) (*AppMetrics, error) {
	metrics := &AppMetrics{}

	var err error

	// HTTP metrics
	metrics.HTTPRequestsTotal, err = meter.Int64Counter(
		"http_requests_total",
		metric.WithDescription("Total number of HTTP requests"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	metrics.HTTPRequestDuration, err = meter.Float64Histogram(
		"http_request_duration_seconds",
		metric.WithDescription("Duration of HTTP requests"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10),
	)
	if err != nil {
		return nil, err
	}

	metrics.HTTPRequestsActive, err = meter.Int64UpDownCounter(
		"http_requests_active",
		metric.WithDescription("Number of active HTTP requests"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	metrics.DBConnectionsOpen, err = meter.Int64Gauge(
		"db_connections_open",
		metric.WithDescription("Total number of open database connections"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	metrics.DBConnectionsInUse, err = meter.Int64Gauge(
		"db_connections_in_use",
		metric.WithDescription("Number of database connections currently in use"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	metrics.DBConnectionsIdle, err = meter.Int64Gauge(
		"db_connections_idle",
		metric.WithDescription("Number of idle database connections"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	metrics.DBQueriesTotal, err = meter.Int64Counter(
		"db_queries_total",
		metric.WithDescription("Total number of database queries"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	metrics.DBQueryDuration, err = meter.Float64Histogram(
		"db_query_duration_seconds",
		metric.WithDescription("Duration of database queries"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5),
	)
	if err != nil {
		return nil, err
	}

	// Auth metrics
	metrics.AuthAttemptsTotal, err = meter.Int64Counter(
		"auth_attempts_total",
		metric.WithDescription("Total number of authentication attempts"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	metrics.AuthFailuresTotal, err = meter.Int64Counter(
		"auth_failures_total",
		metric.WithDescription("Total number of authentication failures"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	// Business metrics
	metrics.UsersTotal, err = meter.Int64UpDownCounter(
		"users_total",
		metric.WithDescription("Total number of users"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	metrics.SimplesTotal, err = meter.Int64UpDownCounter(
		"simples_total",
		metric.WithDescription("Total number of simple entities"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

// HTTP metrics methods
func (m *AppMetrics) RecordHTTPRequest(ctx context.Context, method, path, status string, duration float64) {
	attrs := []attribute.KeyValue{
		attribute.String("method", method),
		attribute.String("path", path),
		attribute.String("status", status),
	}

	m.HTTPRequestsTotal.Add(ctx, 1, metric.WithAttributes(attrs...))
	m.HTTPRequestDuration.Record(ctx, duration, metric.WithAttributes(attrs...))
}

func (m *AppMetrics) RecordActiveHTTPRequest(ctx context.Context, delta int64) {
	m.HTTPRequestsActive.Add(ctx, delta)
}

// Database metrics methods
func (m *AppMetrics) RecordDBQuery(ctx context.Context, operation string, duration float64) {
	attrs := []attribute.KeyValue{
		attribute.String("operation", operation),
	}

	m.DBQueriesTotal.Add(ctx, 1, metric.WithAttributes(attrs...))
	m.DBQueryDuration.Record(ctx, duration, metric.WithAttributes(attrs...))
}

func (m *AppMetrics) RecordDBConnectionStats(ctx context.Context, stats sql.DBStats) {
	m.DBConnectionsOpen.Record(ctx, int64(stats.OpenConnections))
	m.DBConnectionsInUse.Record(ctx, int64(stats.InUse))
	m.DBConnectionsIdle.Record(ctx, int64(stats.Idle))
}

// Auth metrics methods
func (m *AppMetrics) RecordAuthAttempt(ctx context.Context, success bool, method string) {
	attrs := []attribute.KeyValue{
		attribute.String("method", method),
	}

	m.AuthAttemptsTotal.Add(ctx, 1, metric.WithAttributes(attrs...))

	if !success {
		m.AuthFailuresTotal.Add(ctx, 1, metric.WithAttributes(attrs...))
	}
}

// Business metrics methods
func (m *AppMetrics) UpdateUserCount(ctx context.Context, count int64) {
	m.UsersTotal.Add(ctx, count)
}

func (m *AppMetrics) UpdateSimpleCount(ctx context.Context, count int64) {
	m.SimplesTotal.Add(ctx, count)
}
