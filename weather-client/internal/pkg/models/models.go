package models

type WeatherResponse struct {
	CurrentWeather struct {
		Temperature   float64 `json:"temperature"`
		Windspeed     float64 `json:"windspeed"`
		WindDirection float64 `json:"winddirection"`
		IsDay         int     `json:"is_day"`
		WeatherCode   int     `json:"weathercode"`
	} `json:"current_weather"`
}
