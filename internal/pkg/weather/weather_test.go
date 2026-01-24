package weather

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"weather-ingestion-svc/internal/pkg/models"
)

// helper to restore httpGet
func resetHTTPGet() {
	httpGet = http.Get
}

// ---------- MOCK RESPONSES ----------

func mockSuccessResponse(_ string) (*http.Response, error) {
	body := `{
		"current_weather": {
			"temperature": 25.5,
			"windspeed": 12.3,
			"winddirection": 180,
			"is_day": 1,
			"weathercode": 0
		}
	}`
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(body)),
	}, nil
}

func mockFailureResponse(_ string) (*http.Response, error) {
	return nil, errors.New("network error")
}

// ---------- TESTS ----------

func TestFetchWeatherAPIResponse_Success(t *testing.T) {
	httpGet = mockSuccessResponse
	defer resetHTTPGet()

	resp, err := FetchWeatherAPIResponse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.CurrentWeather.Temperature != 25.5 {
		t.Fatalf("unexpected temperature value")
	}
}

func TestFetchWeatherAPIResponse_Error(t *testing.T) {
	httpGet = mockFailureResponse
	defer resetHTTPGet()

	_, err := FetchWeatherAPIResponse()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestScrapeWeatherValues(t *testing.T) {
	weather := &models.WeatherResponse{}
	weather.CurrentWeather.Temperature = 30.0
	weather.CurrentWeather.Windspeed = 15.0

	ScrapeWeatherValues(weather)
}

func TestRunWeatherOnce_Success(t *testing.T) {
	httpGet = mockSuccessResponse
	defer resetHTTPGet()

	err := RunWeatherOnce()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunWeatherOnce_Error(t *testing.T) {
	httpGet = mockFailureResponse
	defer resetHTTPGet()

	err := RunWeatherOnce()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
