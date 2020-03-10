module github.com/gophernment/api

go 1.14

replace github.com/gophernment/api => ./..

require (
	github.com/go-redis/redis/v7 v7.2.0
	github.com/gorilla/mux v1.7.4
	github.com/prometheus/client_golang v1.5.0
	github.com/spf13/viper v1.6.2
	gitlab.com/pallat/api v0.0.0-20200303111155-a15b8016469f
	go.opentelemetry.io/otel v0.2.3
	go.opentelemetry.io/otel/example/grpc v0.2.3
	go.opentelemetry.io/otel/exporter/trace/jaeger v0.2.1
	go.uber.org/zap v1.14.0
	google.golang.org/grpc v1.27.1
)
