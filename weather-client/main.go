package main

import (
	"fmt"
	"net/http"

	"weather-client/internal/pkg/constants"
	"weather-client/internal/pkg/logger"
	"weather-client/internal/pkg/observability"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	logger.Init()
	prometheus.MustRegister(observability.TemperatureGauge)
	prometheus.MustRegister(observability.WindspeedGauge)
	prometheus.MustRegister(observability.APIUpGauge)

	go observability.StartWeatherPolling()

	http.Handle("/metrics", promhttp.Handler())

	fmt.Println("Metrics available at http://localhost" + constants.ServerPort + "/metrics")

	if err := http.ListenAndServe(constants.ServerPort, nil); err != nil {
		fmt.Println("Server failed:", err)
	}
}
