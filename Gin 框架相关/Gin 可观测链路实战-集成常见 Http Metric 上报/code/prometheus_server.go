package main

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"log"
	"net/http"
)

func serveMetrics(prometheusPort int64) {
	http.Handle("/metrics", promhttp.Handler())
	if prometheusPort == 0 {
		prometheusPort = 2223
	}
	addr := fmt.Sprintf(":%d", prometheusPort)
	log.Printf("serving metrics at %s", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Printf("error serving http: %v", err)
		panic(err)
	}
}

func initMetrics(prometheusPort int64, serviceName string) {
	metricExporter, err := prometheus.New()
	if err != nil {
		panic(err)
	}

	res, err := resource.New(context.Background(),
		resource.WithAttributes(semconv.ServiceNameKey.String(serviceName)),
		resource.WithSchemaURL(semconv.SchemaURL),
	)
	if err != nil {
		panic(err)
	}

	meterProvider := metric.NewMeterProvider(metric.WithReader(metricExporter), metric.WithResource(res))
	otel.SetMeterProvider(meterProvider)
	go serveMetrics(prometheusPort)

}
