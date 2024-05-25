package dto

type TemperatureResponseDto struct {
	City       string  `json:"city"`
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kevin      float64 `json:"temp_K"`
}
