package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gophernment/api/logs"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporter/trace/jaeger"
	"go.opentelemetry.io/otel/plugin/othttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

const (
	apiAddressKey = "api.address"
	appENVKey     = "app.env"
)

var (
	buildcommit = "development"
	buildtime   = ""
	version     = "no version"
)

func init() {
	initConfig()
	initTrace()
	logs.InitLogger(viper.GetString(appENVKey) == "dev")
	recordMetrics()
}

func main() {
	defer logs.Sync()

	r := mux.NewRouter()
	r.Use(mux.CORSMethodMiddleware(r))

	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/version", versionHandler)
	r.HandleFunc("/liveness", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })
	r.HandleFunc("/readiness", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	srv := &http.Server{
		Handler: othttp.NewHandler(r, hostname,
			othttp.WithMessageEvents(othttp.ReadEvents, othttp.WriteEvents),
		),
		Addr:         viper.GetString(apiAddressKey),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf(
		"Starting the service...\ncommit: %s, build time: %s, release: %s\n",
		buildcommit, buildtime, version,
	)
	fmt.Printf("serve on %s\n", ":"+viper.GetString(apiAddressKey))
	go func() {
		fmt.Println(srv.ListenAndServe())
	}()

	shutdown(srv)
}

func shutdown(srv *http.Server) {
	sigterm := make(chan os.Signal)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	<-sigterm

	if err := srv.Shutdown(context.Background()); err != nil {
		fmt.Println(err.Error())
	}
}

func recordMetrics() {
	opsProcessed := promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})

	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

func initTrace() {
	exporter, err := jaeger.NewExporter(
		jaeger.WithCollectorEndpoint("http://localhost:14268/api/traces"),
		jaeger.WithProcess(
			jaeger.Process{
				ServiceName: "github.com/pallat/basic",
			},
		),
	)
	if err != nil {
		log.Fatalf("have some errors while creating stdout exporter: %v", err)
	}

	provider, err := sdktrace.NewProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithSyncer(exporter),
	)
	if err != nil {
		log.Fatalf("have some problems while creating provider: %v", err)
	}
	global.SetTraceProvider(provider)
}

func initConfig() {
	viper.SetDefault(apiAddressKey, "0.0.0.0:8080")
	viper.SetDefault(appENVKey, "dev")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("content-type", "application/json")

	resp := map[string]string{
		"buildTime": buildtime,
		"commit":    buildcommit,
		"version":   version,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)
}
