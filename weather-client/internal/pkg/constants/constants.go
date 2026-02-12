package constants

const (
	// Prometheus server port
	ServerPort = ":2112"

	// Weather API base (without lat/long)
	WeatherAPIBaseURL = "https://api.open-meteo.com/v1/forecast"

	// PostgreSQL connection string
	DatabaseURL = "postgres://weather_user:weather123@localhost:5433/weather_db?sslmode=disable"
)
