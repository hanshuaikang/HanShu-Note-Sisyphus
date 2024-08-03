package main

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"time"
)

type httpMetricsRecorder struct {
	requestsCounter metric.Int64UpDownCounter
	totalDuration   metric.Int64Histogram
}

func NewHttpMetricsRecorder(instrumentationName, metricsPrefix string) httpMetricsRecorder {
	metricName := func(metricName string) string {
		if len(metricsPrefix) > 0 {
			return metricsPrefix + "." + metricName
		}
		return metricName
	}
	meter := otel.Meter(instrumentationName, metric.WithInstrumentationVersion("semver:1.0.0"))
	requestsCounter, _ := meter.Int64UpDownCounter(metricName("http.server.request_count"), metric.WithDescription("Number of Requests"), metric.WithUnit("Count"))
	totalDuration, _ := meter.Int64Histogram(metricName("http.server.duration"), metric.WithDescription("Time Taken by request"), metric.WithUnit("Milliseconds"))

	return httpMetricsRecorder{
		requestsCounter: requestsCounter,
		totalDuration:   totalDuration,
	}
}

// AddRequests 请求开始的时候 调用这个函数为requestsCounter 这个计数器 + 1
func (r *httpMetricsRecorder) AddRequests(ctx context.Context, attributes []attribute.KeyValue) {
	r.requestsCounter.Add(ctx, 1, metric.WithAttributes(attributes...))
}

// ObserveHTTPRequestDuration 这里接受一个参数,表示请求的持续时间
func (r *httpMetricsRecorder) ObserveHTTPRequestDuration(ctx context.Context, duration time.Duration, attributes []attribute.KeyValue) {
	r.totalDuration.Record(ctx, int64(duration/time.Millisecond), metric.WithAttributes(attributes...))
}
