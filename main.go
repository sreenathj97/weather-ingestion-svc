package main

import (
	"fmt"
	"net/http"

	"weather-ingestion-svc/internal/pkg/constants"
	"weather-ingestion-svc/internal/pkg/observability"
	"weather-ingestion-svc/internal/pkg/weather"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	prometheus.MustRegister(observability.TemperatureGauge)
	prometheus.MustRegister(observability.WindspeedGauge)
	prometheus.MustRegister(observability.APIUpGauge)

	go weather.StartWeatherPolling()

	http.Handle("/metrics", promhttp.Handler())

	fmt.Println("Metrics available at http://localhost" + constants.ServerPort + "/metrics")

	if err := http.ListenAndServe(constants.ServerPort, nil); err != nil {
		fmt.Println("Server failed:", err)
	}
}
