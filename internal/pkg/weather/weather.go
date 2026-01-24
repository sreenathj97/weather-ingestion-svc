package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"weather-ingestion-svc/internal/pkg/constants"
	"weather-ingestion-svc/internal/pkg/logger"
	"weather-ingestion-svc/internal/pkg/models"
	"weather-ingestion-svc/internal/pkg/observability"
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
	observability.TemperatureGauge.Set(weather.CurrentWeather.Temperature)
	observability.WindspeedGauge.Set(weather.CurrentWeather.Windspeed)
	observability.APIUpGauge.Set(1)
}

// RunWeatherOnce executes one polling cycle (TESTABLE)
func RunWeatherOnce() error {
	weather, err := FetchWeatherAPIResponse()
	if err != nil {
		observability.APIUpGauge.Set(0)
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
