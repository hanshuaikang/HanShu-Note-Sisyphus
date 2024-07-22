大家好，我是韩数，本身自己是做 Devops 开发的，之前也研究过一段时间的 Influxdb , 也算是接触过一点点可观测相关的技术，这几年 Opentelemetry 很流行，于是自己就想能不能能不能写一个小工具快速集成 gin上报一些基础的 http 性能指标到 Prometheus 去呢，去 Github 上找了一圈, 发现很多库都已经不更新了, 而且也都是基于 Prometheus Client 去实现的，而 Opentelemetry 的文档, 懂得都懂, 网上很难找到比较细致的教程和 demo。一切都得自己慢慢踩坑, 所以刚好趁着折腾 Opentelemetry 的契机，何不基于 Opentelemetry 的 golang sdk 实现一个 `gin-prometheus` 项目呢, 说干就干。 

当然熟悉我的人都知道(笑死，掘金上怎么可能有人会熟悉我), 我肯定不会就此善罢甘休的，之后我会更新两篇文章，分享我自己使用 gin 去接入 Opentelemetry Trace 和 Metrics 的心路历程，当然代码也会附带分享出来, 希望能帮助对那些想要了解 Opentelemetry 的人提供一些微弱的帮助，为什么现在不更呢，因为还没写。

![](./img/img.png)

## gin-prometheus 介绍

仓库地址: https://github.com/hanshuaikang/gin-prometheus

gin-prometheus 是一个基于 Opentelemetry 实现的工具，用于简化开发者的配置，gin 项目只需要一行代码，便可以集成常见的Http性能指标。同时 gin-prometheus 提供了丰富的可选的配置，开发者可以选择关闭或者开启某些指标的采集。

🚀 Features:
- 基于 Opentelemetry 实现, 可以作为学习 Opentelemetry 的 案例
- 集成常用的 http 延时，请求数，请求大小等指标。
- 集成 cpu 使用率 & 内存使用率 系统指标。
- 支持自定义的 Recorder, 替换原有的 otel Recorder 实现。
- 丰富的自定义配置: 可以选择只保留部分指标&支持自定义指标前缀&全局 label &prometheus端口...

🚀 后续计划:
- 支持 prometheus push & console log 的方式暴露 metrics。
- 支持 更多的 http 指标暴露
- 等


### 🎉 Metrics

Details about exposed Prometheus metrics.

| Name                                   | Type | Exposed Information           |
|----------------------------------------| ---- |-------------------------------|
| http_server_active_requests                        | Counter    | Number of requests inflight.  |
| http_server_duration                      | Histogram   | Time Taken by request         |
| http_server_request_total              | Counter | Number of Requests.           |
| http_server_request_content_length         | Histogram  | HTTP request sizes in bytes.  |
| http_server_response_content_length       | Histogram   | HTTP response sizes in bytes. |
| system_cpu_usage                         | Counter  | CPU Usage                     |
| system_memory_usage                        | Counter    | Memory Usage                  |


## Usage

接入 gin-prometheus 非常简单，只需要 安装 gin-prometheus 并且注册 gin-prometheus 提供的中间件即可。

```bash
go get -u github.com/hanshuaikang/gin-prometheus@v0.0.1
```

在 gin 的 engine 中注入 gin-prometheus 中间件:

```golang
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/hanshuaikang/gin-prometheus"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"net/http"
)

const serviceName = "gin-prometheus-demo"

func main() {
	fmt.Println("initializing")
	r := gin.New()
	r.Use(ginprometheus.Middleware())
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}
```

### 自定义更多的配置:

gin-prometheus 提供了不算丰富的可选配置, 如果想要自定义更多的配置，可以参考如下：
```golang

func main() {
	fmt.Println("initializing")

	r := gin.New()
	globalAttributes := []attribute.KeyValue{
		semconv.K8SPodName("pod-1"),
		semconv.K8SNamespaceName("test"),
		semconv.ServiceName(serviceName),
	}
	r.Use(ginprometheus.Middleware(
		## 自定义指标属性配置，你可以选择注入某些自定义的 label 到 metric
		ginprometheus.WithAttributes(func(route string, request *http.Request) []attribute.KeyValue {
			attrs := []attribute.KeyValue{
				semconv.HTTPMethodKey.String(request.Method),
			}
			if route != "" {
				attrs = append(attrs, semconv.HTTPRouteKey.String(route))
			}
			return attrs
		}),
           ## 注册全局 label, 所有的 metric 都会携带这些 label
		ginprometheus.WithGlobalAttributes(globalAttributes),
           ## 自定义服务信息      
		ginprometheus.WithService(serviceName, "v0.0.1"),
           ## infra 自定义指标前缀
		ginprometheus.WithMetricPrefix("infra"),
           ## 自定义 Prometheus 端口
		ginprometheus.WithPrometheusPort(4433),
           ## 关闭 cpu 和 内存使用率的 采集
		ginprometheus.WithSystemMetricDisabled(),
           ## ....
	))
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}

```







