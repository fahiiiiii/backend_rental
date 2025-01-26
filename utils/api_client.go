package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/beego/beego/v2/core/config"
	"backend_rental/models"
)

type ApiClient struct {
	BaseURL string
	Headers map[string]string
}

func NewApiClient() *ApiClient {
	// Use Beego's configuration to get RapidAPI key
	conf, err := config.NewConfig("ini", "conf/app.conf")
	if err != nil {
		fmt.Printf("Error reading configuration: %v\n", err)
		return nil
	}

	// Get the RapidAPI key, case-insensitive
	apiKey, err := conf.String("RAPIDAPI_KEY")
	if err != nil || apiKey == "" {
		fmt.Printf("Warning: RapidAPI key not found or is empty\n")
		apiKey = "default_key_if_needed"
	}

	return &ApiClient{
		BaseURL: "https://booking-com18.p.rapidapi.com/stays/auto-complete",
		Headers: map[string]string{
			"x-rapidapi-host": "booking-com18.p.rapidapi.com",
			"x-rapidapi-key":  strings.TrimSpace(apiKey),
		},
	}
}

func (c *ApiClient) FetchCityData(query string) (*models.ApiResponse, error) {
	url := fmt.Sprintf("%s?query=%s", c.BaseURL, query)
	req, _ := http.NewRequest("GET", url, nil)

	// Set headers
	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse response
	body, _ := ioutil.ReadAll(resp.Body)
	var apiResponse models.ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	return &apiResponse, nil
}
