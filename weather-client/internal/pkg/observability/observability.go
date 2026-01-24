package observability

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"weather-client/internal/pkg/constants"
	"weather-client/internal/pkg/logger"
	"weather-client/internal/pkg/models"

	"github.com/prometheus/client_golang/prometheus"
)

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

// httpGet is a variable so it can be mocked in tests
var httpGet = http.Get

// FetchWeatherAPIResponse calls the weather API and parses response
func FetchWeatherAPIResponse() (*models.WeatherResponse, error) {
	resp, err := httpGet(constants.WeatherAPIURL)
	if err != nil {
		return nil, fmt.Errorf("api call failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 response: %s", resp.Status)
	}

	var weather models.WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return nil, fmt.Errorf("json decode failed: %v", err)
	}

	return &weather, nil
}

// ScrapeWeatherValues updates Prometheus metrics
func ScrapeWeatherValues(weather *models.WeatherResponse) {
	TemperatureGauge.Set(weather.CurrentWeather.Temperature)
	WindspeedGauge.Set(weather.CurrentWeather.Windspeed)
	APIUpGauge.Set(1)
}

// RunWeatherOnce executes one polling cycle (TESTABLE)
func RunWeatherOnce() error {
	weather, err := FetchWeatherAPIResponse()
	if err != nil {
		APIUpGauge.Set(0)
		return err
	}

	ScrapeWeatherValues(weather)
	return nil
}

// StartWeatherPolling runs continuously (NOT unit tested)
func StartWeatherPolling() {
	for {
		if err := RunWeatherOnce(); err != nil {
			logger.Logger.Error(
				"Weather fetch failed",
				"error", err,
			)

		} else {
			logger.Logger.Info(
				"Weather updated successfully",
			)

		}
		time.Sleep(5 * time.Second)
	}
}
