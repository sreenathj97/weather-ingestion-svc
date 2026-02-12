package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"weather-client/internal/pkg/constants"
	"weather-client/internal/pkg/logger"
	"weather-client/internal/pkg/observability"
	"weather-client/internal/pkg/repository"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const pollInterval = 5 * time.Second

func main() {
	logger.Init()

	db, err := sql.Open("postgres", constants.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect database", "error", err)
		return
	}
	defer db.Close()

	repo := repository.NewCityRepository(db)

	prometheus.MustRegister(
		observability.TemperatureGauge,
		observability.WindspeedGauge,
	)

	observabilityClient := observability.NewObservabilityClient(
		constants.WeatherAPIBaseURL,
		http.DefaultClient,
		repo,
	)

	go observabilityClient.WeatherMetricsWorkflow(pollInterval)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	httpServer := &http.Server{
		Addr:    constants.ServerPort,
		Handler: mux,
	}

	slog.Info("starting metrics server", "addr", httpServer.Addr)

	if err := httpServer.ListenAndServe(); err != nil {
		slog.Error("server stopped", "error", err)
	}
}
