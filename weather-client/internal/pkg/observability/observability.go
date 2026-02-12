package observability

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"weather-client/internal/pkg/logger"
	"weather-client/internal/pkg/models"
	"weather-client/internal/pkg/repository"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	TemperatureGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "weather_temperature_celsius",
			Help: "Current temperature in Celsius",
		},
		[]string{"city"},
	)

	WindspeedGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "weather_windspeed_kmh",
			Help: "Current wind speed in km/h",
		},
		[]string{"city"},
	)
)

type ObservabilityInterface interface {
	WeatherMetricsWorkflow(interval time.Duration)
}

type httpGetter interface {
	Get(url string) (*http.Response, error)
}

type observabilityClient struct {
	weatherApiBaseUrl string
	http              httpGetter
	dbRepo            repository.CityRepository
}

func NewObservabilityClient(
	apiBaseUrl string,
	httpClient httpGetter,
	dbRepo repository.CityRepository,
) ObservabilityInterface {

	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &observabilityClient{
		weatherApiBaseUrl: apiBaseUrl,
		http:              httpClient,
		dbRepo:            dbRepo,
	}
}

func (c *observabilityClient) WeatherMetricsWorkflow(interval time.Duration) {

	logger.Logger.Info("weather metrics workflow started", "interval", interval)

	for {
		ctx := context.Background()

		cities, err := c.dbRepo.GetAllCities(ctx)
		if err != nil {
			logger.Logger.Error("failed to fetch cities", "error", err)
			time.Sleep(interval)
			continue
		}

		var wg sync.WaitGroup

		for _, city := range cities {
			wg.Add(1)

			go func(city models.City) {
				defer wg.Done()
				c.fetchAndEmit(city)
			}(city)
		}

		wg.Wait()
		time.Sleep(interval)
	}
}

func (c *observabilityClient) fetchAndEmit(city models.City) {

	url := fmt.Sprintf(
		"%s?latitude=%f&longitude=%f&current_weather=true",
		c.weatherApiBaseUrl,
		city.Latitude,
		city.Longitude,
	)

	resp, err := c.http.Get(url)
	if err != nil {
		logger.Logger.Warn("weather api failed", "city", city.CityName, "error", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Logger.Warn("non-200 response", "city", city.CityName)
		return
	}

	var weather models.WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		logger.Logger.Warn("json decode failed", "city", city.CityName)
		return
	}

	TemperatureGauge.WithLabelValues(city.CityName).Set(weather.CurrentWeather.Temperature)
	WindspeedGauge.WithLabelValues(city.CityName).Set(weather.CurrentWeather.Windspeed)

	logger.Logger.Debug("metrics emitted", "city", city.CityName)
}
