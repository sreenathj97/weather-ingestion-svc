package observability

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockHTTPClient struct {
	mock.Mock
}

func (m *mockHTTPClient) Get(url string) (*http.Response, error) {
	args := m.Called(url)
	resp, _ := args.Get(0).(*http.Response)
	return resp, args.Error(1)
}

func TestGetWeatherMetricsSuccess(t *testing.T) {
	t.Parallel()

	const apiURL = "https://example.com/weather"
	body := `{"current_weather":{"temperature":21.5,"windspeed":9.2,"winddirection":120,"is_day":1,"weathercode":3}}`

	mockClient := &mockHTTPClient{}
	mockClient.On("Get", apiURL).Return(&http.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}, nil)

	client := NewObservabilityClient(apiURL, mockClient)
	weatherApiResponse, err := client.GetWeatherMetrics()

	require.NoError(t, err)
	require.NotNil(t, weatherApiResponse)
	require.Equal(t, 21.5, weatherApiResponse.CurrentWeather.Temperature)
	require.Equal(t, 9.2, weatherApiResponse.CurrentWeather.Windspeed)
	require.Equal(t, 120.0, weatherApiResponse.CurrentWeather.WindDirection)
	require.Equal(t, 1, weatherApiResponse.CurrentWeather.IsDay)
	require.Equal(t, 3, weatherApiResponse.CurrentWeather.WeatherCode)

	mockClient.AssertExpectations(t)
}

func TestGetWeatherMetricsHTTPError(t *testing.T) {
	t.Parallel()

	const apiURL = "https://example.com/weather"

	mockClient := &mockHTTPClient{}
	mockClient.On("Get", apiURL).Return((*http.Response)(nil), errors.New("network down"))

	client := NewObservabilityClient(apiURL, mockClient)
	weatherApiResponse, err := client.GetWeatherMetrics()

	require.Error(t, err)
	require.Nil(t, weatherApiResponse)
	require.Contains(t, err.Error(), "api call failed")

	mockClient.AssertExpectations(t)
}

func TestGetWeatherMetricsNon200(t *testing.T) {
	t.Parallel()

	const apiURL = "https://example.com/weather"

	mockClient := &mockHTTPClient{}
	mockClient.On("Get", apiURL).Return(&http.Response{
		StatusCode: http.StatusInternalServerError,
		Status:     "500 Internal Server Error",
		Body:       io.NopCloser(strings.NewReader(`{}`)),
		Header:     make(http.Header),
	}, nil)

	client := NewObservabilityClient(apiURL, mockClient)
	weatherApiResponse, err := client.GetWeatherMetrics()

	require.Error(t, err)
	require.Nil(t, weatherApiResponse)
	require.Contains(t, err.Error(), "non-200 response")

	mockClient.AssertExpectations(t)
}

func TestGetWeatherMetricsInvalidJSON(t *testing.T) {
	t.Parallel()

	const apiURL = "https://example.com/weather"

	mockClient := &mockHTTPClient{}
	mockClient.On("Get", apiURL).Return(&http.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body:       io.NopCloser(strings.NewReader("{")),
		Header:     make(http.Header),
	}, nil)

	client := NewObservabilityClient(apiURL, mockClient)
	weatherApiResponse, err := client.GetWeatherMetrics()

	require.Error(t, err)
	require.Nil(t, weatherApiResponse)
	require.Contains(t, err.Error(), "json decode failed")

	mockClient.AssertExpectations(t)
}
