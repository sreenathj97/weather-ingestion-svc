package observability

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"weather-client/internal/pkg/models"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	TemperatureGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "weather_temperature_celsius",
		Help: "Current temperature in Celsius",
	})

	WindspeedGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "weather_windspeed_kmh",
		Help: "Current wind speed in km/h",
	})

	logger = slog.Default()
)

type ObservabilityInterface interface {
	WeatherMetricsWorkflow(interval time.Duration)
	GetWeatherMetrics() (*models.WeatherResponse, error)
}

type httpGetter interface {
	Get(url string) (*http.Response, error)
}

type observabilityClient struct {
	weatherApiUrl string
	http          httpGetter
}

func NewObservabilityClient(weatherApiUrl string, httpClient httpGetter) ObservabilityInterface {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &observabilityClient{
		weatherApiUrl: weatherApiUrl,
		http:          httpClient,
	}
}

func (c *observabilityClient) GetWeatherMetrics() (*models.WeatherResponse, error) {

	weatherApiResponse, err := c.http.Get(c.weatherApiUrl)
	if err != nil {
		logger.Warn("weather api request failed", "error", err)
		return nil, fmt.Errorf("api call failed: %w", err)
	}
	defer weatherApiResponse.Body.Close()

	if weatherApiResponse.StatusCode != http.StatusOK {
		logger.Warn(
			"weather api returned non-200 response",
			"status", weatherApiResponse.Status,
		)
		return nil, fmt.Errorf("non-200 response: %s", weatherApiResponse.Status)
	}

	var weatherDetails models.WeatherResponse
	if err := json.NewDecoder(weatherApiResponse.Body).Decode(&weatherDetails); err != nil {
		logger.Warn("failed to decode weather api response", "error", err)
		return nil, fmt.Errorf("json decode failed: %w", err)
	}

	return &weatherDetails, nil
}

func EmitWeatherMetrics(weatherMetrics *models.WeatherResponse) {
	if weatherMetrics == nil {
		logger.Warn("skipping metric emission: nil weather response")
		return
	}

	TemperatureGauge.Set(weatherMetrics.CurrentWeather.Temperature)
	WindspeedGauge.Set(weatherMetrics.CurrentWeather.Windspeed)

	logger.Debug(
		"weather metrics emitted",
		"temperature_c", weatherMetrics.CurrentWeather.Temperature,
		"windspeed_kmh", weatherMetrics.CurrentWeather.Windspeed,
	)
}

func (c *observabilityClient) WeatherMetricsWorkflow(interval time.Duration) {

	logger.Info("weather metrics workflow started", "interval", interval)

	for {
		weatherMetrics, err := c.GetWeatherMetrics()
		if err == nil {
			EmitWeatherMetrics(weatherMetrics)
		}
		time.Sleep(interval)
	}
}
