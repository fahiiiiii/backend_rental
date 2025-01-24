package utils

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "time"
    "backend_rental/models"
)

func (c *ApiClient) FetchPropertiesForCity(locationId string, checkIn, checkOut string) (*models.PropertyResponse, error) {
    url := fmt.Sprintf(
        "https://booking-com18.p.rapidapi.com/stays/search?locationId=%s&checkinDate=%s&checkoutDate=%s&units=metric&temperature=c", 
        locationId, 
        checkIn, 
        checkOut,
    )
    
    req, _ := http.NewRequest("GET", url, nil)

    // Set headers
    for key, value := range c.Headers {
        req.Header.Set(key, value)
    }

    // Send request
    client := &http.Client{
        Timeout: 30 * time.Second,
    }
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // Parse response
    body, _ := ioutil.ReadAll(resp.Body)
    var propertyResponse models.PropertyResponse
    err = json.Unmarshal(body, &propertyResponse)
    if err != nil {
        return nil, err
    }

    return &propertyResponse, nil
}