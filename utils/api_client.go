package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"backend_rental/models"
)

type ApiClient struct {
	BaseURL string
	Headers map[string]string
}

func NewApiClient() *ApiClient {
	return &ApiClient{
		BaseURL: "https://booking-com18.p.rapidapi.com/stays/auto-complete",
		Headers: map[string]string{
			"x-rapidapi-host": "booking-com18.p.rapidapi.com",
			"x-rapidapi-key":  "d331a74390msh5ef004857f77ed7p133e3ejsn0ae0b7a6a96c",
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
