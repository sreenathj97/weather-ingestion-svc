package observability

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"weather-client/internal/pkg/models"

	"github.com/stretchr/testify/mock"
)

//
// ---------- MOCK HTTP CLIENT ----------
//

type mockHTTPClient struct {
	mock.Mock
}

func (m *mockHTTPClient) Get(url string) (*http.Response, error) {
	args := m.Called(url)
	resp, _ := args.Get(0).(*http.Response)
	return resp, args.Error(1)
}

//
// ---------- MOCK REPOSITORY ----------
//

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) GetAllCities(ctx context.Context) ([]models.City, error) {
	args := m.Called(ctx)
	cities, _ := args.Get(0).([]models.City)
	return cities, args.Error(1)
}

//
// ---------- TEST SUCCESS CASE ----------
//

func TestWeatherWorkflowSuccess(t *testing.T) {
	t.Parallel()

	apiBaseURL := "https://example.com/weather"

	mockClient := new(mockHTTPClient)
	mockRepository := new(mockRepo)

	testCities := []models.City{
		{
			CityName:  "TestCity",
			Latitude:  10.0,
			Longitude: 20.0,
		},
	}

	mockRepository.
		On("GetAllCities", mock.Anything).
		Return(testCities, nil)

	body := `{"current_weather":{"temperature":25.5,"windspeed":5.2}}`

	mockClient.
		On("Get", mock.Anything).
		Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil)

	client := NewObservabilityClient(apiBaseURL, mockClient, mockRepository)

	// Run only one iteration safely
	go client.WeatherMetricsWorkflow(100 * time.Millisecond)

	time.Sleep(200 * time.Millisecond)

	mockClient.AssertExpectations(t)
	mockRepository.AssertExpectations(t)
}

//
// ---------- TEST DB ERROR ----------
//

func TestWeatherWorkflowDBError(t *testing.T) {
	t.Parallel()

	apiBaseURL := "https://example.com/weather"

	mockClient := new(mockHTTPClient)
	mockRepository := new(mockRepo)

	mockRepository.
		On("GetAllCities", mock.Anything).
		Return(nil, context.DeadlineExceeded)

	client := NewObservabilityClient(apiBaseURL, mockClient, mockRepository)

	go client.WeatherMetricsWorkflow(100 * time.Millisecond)

	time.Sleep(200 * time.Millisecond)

	mockRepository.AssertExpectations(t)
}

//
// ---------- TEST NON-200 RESPONSE ----------
//

func TestWeatherWorkflowNon200(t *testing.T) {
	t.Parallel()

	apiBaseURL := "https://example.com/weather"

	mockClient := new(mockHTTPClient)
	mockRepository := new(mockRepo)

	testCities := []models.City{
		{
			CityName:  "TestCity",
			Latitude:  10.0,
			Longitude: 20.0,
		},
	}

	mockRepository.
		On("GetAllCities", mock.Anything).
		Return(testCities, nil)

	mockClient.
		On("Get", mock.Anything).
		Return(&http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(strings.NewReader(`{}`)),
			Header:     make(http.Header),
		}, nil)

	client := NewObservabilityClient(apiBaseURL, mockClient, mockRepository)

	go client.WeatherMetricsWorkflow(100 * time.Millisecond)

	time.Sleep(200 * time.Millisecond)

	mockClient.AssertExpectations(t)
	mockRepository.AssertExpectations(t)
}
