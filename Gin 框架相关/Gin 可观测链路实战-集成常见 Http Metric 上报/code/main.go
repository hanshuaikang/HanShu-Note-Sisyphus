package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Listen and Server in 0.0.0.0:8080
	r := gin.Default()
	r.Use(HttpMetricMiddleware())
	initMetrics(2233, "gin_metric_name")
	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	r.Run(":8080")
}
