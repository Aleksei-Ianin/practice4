package domain

import "time"

type Calculation struct {
	ID        int       `json:"id"`
	Operand1  float64   `json:"operand1"`
	Operand2  float64   `json:"operand2"`
	Operation string    `json:"operation"`   // "add", "subtract", "multiply", "divide", "power"
	Result    float64   `json:"result"`
	CreatedAt time.Time `json:"created_at"`
}