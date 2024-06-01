package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Insurance struct {
	InsuranceData InsuranceData `json:"insurance"`
}

type InsuranceData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Image       string `json:"image"`
}

func GetData(url string) (*Insurance, error) {
	client := &http.Client{Timeout: time.Minute * 5} // Max timeout 5 minutes
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected client response %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var insurance Insurance
	if err := json.Unmarshal(body, &insurance); err != nil {
		return nil, err
	}

	return &insurance, nil
}
