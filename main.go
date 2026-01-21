

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// 1️⃣ STRUCT TO MATCH API RESPONSE
type WeatherResponse struct {
	CurrentWeather struct {
		Temperature   float64 `json:"temperature"`
		Windspeed     float64 `json:"windspeed"`
		WindDirection float64 `json:"winddirection"`
		IsDay         int     `json:"is_day"`
		WeatherCode   int     `json:"weathercode"`
	} `json:"current_weather"`
}

// 2️⃣ PROMETHEUS METRICS (GAUGES)
var (
	temperatureGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "weather_temperature_celsius",
			Help: "Current temperature in Celsius",
		},
	)

	windspeedGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "weather_windspeed_kmh",
			Help: "Current wind speed in km/h",
		},
	)

	apiUpGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "weather_api_up",
			Help: "API status (1 = API working, 0 = API failed)",
		},
	)
)

// 3️⃣ FUNCTION THAT CALLS WEATHER API EVERY 5 SECONDS
func fetchWeather() {
	apiURL := "https://api.open-meteo.com/v1/forecast?latitude=52.52&longitude=13.41&current_weather=true"

	for {
		resp, err := http.Get(apiURL)
		if err != nil {
			fmt.Println("❌ API call failed:", err)
			apiUpGauge.Set(0)
			time.Sleep(5 * time.Second)
			continue
		}

		// Check HTTP status
		if resp.StatusCode != http.StatusOK {
			fmt.Println("❌ Non-200 response:", resp.Status)
			apiUpGauge.Set(0)
			resp.Body.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		var weather WeatherResponse
		err = json.NewDecoder(resp.Body).Decode(&weather)
		resp.Body.Close()

		if err != nil {
			fmt.Println("❌ JSON decode failed:", err)
			apiUpGauge.Set(0)
			time.Sleep(5 * time.Second)
			continue
		}

		// Update Prometheus metrics
		temperatureGauge.Set(weather.CurrentWeather.Temperature)
		windspeedGauge.Set(weather.CurrentWeather.Windspeed)
		apiUpGauge.Set(1)

		fmt.Printf(
			"✅ Weather updated | Temp: %.1f°C | Wind: %.1f km/h\n",
			weather.CurrentWeather.Temperature,
			weather.CurrentWeather.Windspeed,
		)

		time.Sleep(5 * time.Second)
	}
}

// 4️⃣ MAIN FUNCTION (PROGRAM STARTS HERE)
func main() {
	// Register metrics
	prometheus.MustRegister(temperatureGauge)
	prometheus.MustRegister(windspeedGauge)
	prometheus.MustRegister(apiUpGauge)

	// Start weather fetcher in background
	go fetchWeather()

	// Expose /metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	fmt.Println("🚀 Metrics available at http://localhost:2112/metrics")

	// Start HTTP server
	err := http.ListenAndServe(":2112", nil)
	if err != nil {
		fmt.Println("❌ Server failed:", err)
	}
}


