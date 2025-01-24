package utils

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "time"
    "backend_rental/models"
)

func (c *ApiClient) FetchPropertyImages(hotelId string) (*models.PropertyImageResponse, error) {
    url := fmt.Sprintf(
        "https://booking-com18.p.rapidapi.com/stays/get-photos?hotelId=%s", 
        hotelId,
    )
    
    req, _ := http.NewRequest("GET", url, nil)

    for key, value := range c.Headers {
        req.Header.Set(key, value)
    }

    client := &http.Client{
        Timeout: 30 * time.Second,
    }
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)
    var imageResponse models.PropertyImageResponse
    err = json.Unmarshal(body, &imageResponse)
    if err != nil {
        return nil, err
    }

    return &imageResponse, nil
}