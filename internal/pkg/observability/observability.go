package observability

import "github.com/prometheus/client_golang/prometheus"

var (
	TemperatureGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "weather_temperature_celsius",
			Help: "Current temperature in Celsius",
		},
	)

	WindspeedGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "weather_windspeed_kmh",
			Help: "Current wind speed in km/h",
		},
	)

	APIUpGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "weather_api_up",
			Help: "API status (1 = up, 0 = down)",
		},
	)
)
