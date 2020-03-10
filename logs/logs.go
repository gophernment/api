package logs

import (
	"context"
	"net/http"
	"os"

	"go.opentelemetry.io/otel/plugin/grpctrace"
	"go.opentelemetry.io/otel/plugin/httptrace"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

const (
	traceIDKey = "trace-id"
)

var L = zap.NewExample().Sugar()

func Sync() error {
	return L.Sync()
}

func InitLogger(prod bool) {
	cfg := zap.NewDevelopmentConfig()
	if prod {
		cfg = zap.NewProductionConfig()
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	hn, _ := os.Hostname()
	if hn == "" {
		hn = "unknown"
	}

	L = logger.Sugar().With(zap.String("hostname", hn))
}

func NewWithGRPCContext(ctx context.Context) *zap.SugaredLogger {
	requestMetadata, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return L
	}

	metadataCopy := requestMetadata.Copy()
	_, spanCtx := grpctrace.Extract(ctx, &metadataCopy)

	return L.With(traceIDKey, spanCtx.TraceIDString())
}

func NewWithHTTPContext(req *http.Request) *zap.SugaredLogger {
	_, _, span := httptrace.Extract(req.Context(), req)

	return L.With(traceIDKey, span.TraceIDString())
}
