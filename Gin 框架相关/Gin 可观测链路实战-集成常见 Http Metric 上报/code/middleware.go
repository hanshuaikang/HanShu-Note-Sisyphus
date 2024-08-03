package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"time"
)

// HttpMetricMiddleware 请求中间件
func HttpMetricMiddleware() gin.HandlerFunc {
	ctx := context.Background()

	// 初始化记录器
	recorder := NewHttpMetricsRecorder("gin_metric_demo", "")

	return func(ginCtx *gin.Context) {
		// 获取这次请求的完整路径
		route := ginCtx.FullPath()
		if len(route) <= 0 {
			ginCtx.Next()
			return
		}

		// 记录请求开始的时间
		start := time.Now()

		defer func() {
			// 这里我们定义三个 label, 分别为请求方法和请求路径和状态码
			attributes := []attribute.KeyValue{
				semconv.HTTPMethodKey.String(ginCtx.Request.Method),   // 请求方法
				semconv.HTTPRouteKey.String(route),                    // 请求路径
				semconv.HTTPStatusCodeKey.Int(ginCtx.Writer.Status())} // 状态码

			// 请求记录器 + 1
			recorder.AddRequests(ctx, attributes)
			// 记录请求的耗时
			recorder.ObserveHTTPRequestDuration(ctx, time.Since(start), attributes)

		}()

		ginCtx.Next()
	}
}
