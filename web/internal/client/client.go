package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Calculation struct {
	ID        int       `json:"id"`
	Operand1  float64   `json:"operand1"`
	Operand2  float64   `json:"operand2"`
	Operation string    `json:"operation"`
	Result    float64   `json:"result"`
	CreatedAt time.Time `json:"created_at"`
}

type APIClient struct {
	baseURL string
	client  *http.Client
}

func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *APIClient) Calculate(operand1, operand2 float64, operation string) (*Calculation, error) {
	reqBody := map[string]interface{}{
		"operand1":  operand1,
		"operand2":  operand2,
		"operation": operation,
	}
	jsonData, _ := json.Marshal(reqBody)
	resp, err := c.client.Post(c.baseURL+"/api/calculate", "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}
	var calc Calculation
	if err := json.NewDecoder(resp.Body).Decode(&calc); err != nil {
		return nil, err
	}
	return &calc, nil
}

func (c *APIClient) GetAllCalculations() ([]Calculation, error) {
	resp, err := c.client.Get(c.baseURL + "/api/calculations")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var calcs []Calculation
	if err := json.NewDecoder(resp.Body).Decode(&calcs); err != nil {
		return nil, err
	}
	return calcs, nil
}

func (c *APIClient) DeleteCalculation(id int) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/calculations/%d", c.baseURL, id), nil)
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("delete failed: %s", resp.Status)
	}
	return nil
}