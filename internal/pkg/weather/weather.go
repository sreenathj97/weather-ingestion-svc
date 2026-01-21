package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"weather-ingestion-svc/internal/pkg/constants"
	"weather-ingestion-svc/internal/pkg/models"
	"weather-ingestion-svc/internal/pkg/observability"
)

// FetchWeatherAPIResponse handles API call
func FetchWeatherAPIResponse() (*models.WeatherResponse, error) {
	resp, err := http.Get(constants.WeatherAPIURL)
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

// StartWeatherPolling runs the loop
func StartWeatherPolling() {
	for {
		weather, err := FetchWeatherAPIResponse()
		if err != nil {
			observability.APIUpGauge.Set(0)
			fmt.Println(fmt.Sprintf("Weather fetch error: %v", err))
			time.Sleep(5 * time.Second)
			continue
		}

		ScrapeWeatherValues(weather)

		fmt.Println(fmt.Sprintf(
			"Weather updated | Temp: %.1f°C | Wind: %.1f km/h",
			weather.CurrentWeather.Temperature,
			weather.CurrentWeather.Windspeed,
		))

		time.Sleep(5 * time.Second)
	}
}
