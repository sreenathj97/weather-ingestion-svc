package main

import (
	"log/slog"
	"net/http"
	"time"

	"weather-client/internal/pkg/constants"
	"weather-client/internal/pkg/observability"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	pollInterval = 5 * time.Second
)

func main() {
	prometheus.MustRegister(
		observability.TemperatureGauge,
		observability.WindspeedGauge,
	)

	observabilityClient := observability.NewObservabilityClient(constants.WeatherAPIURL, http.DefaultClient)
	go observabilityClient.WeatherMetricsWorkflow(pollInterval)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	httpServer := &http.Server{
		Addr:    constants.ServerPort,
		Handler: mux,
	}

	slog.Info("starting metrics server", "addr", httpServer.Addr)

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("metrics server stopped unexpectedly", "error", err)
	}
}
