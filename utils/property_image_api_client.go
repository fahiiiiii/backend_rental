package utils

import (
    "fmt"
    "net/http"
    "io"
    "encoding/json"
	"os"
)

type ImageApiClient struct {
    Client *http.Client
}

func NewImageApiClient() *ImageApiClient {
    return &ImageApiClient{
        Client: &http.Client{},
    }
}

func (c *ImageApiClient) FetchPropertyImages(hotelID int) (map[string]interface{}, error) {
    url := fmt.Sprintf("https://booking-com18.p.rapidapi.com/stays/get-photos?hotelId=%d", hotelID)
    
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %v", err)
    }
    
    req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
    // IMPORTANT: Replace with environment variable or secure config management
    req.Header.Add("x-rapidapi-key", os.Getenv("RAPIDAPI_KEY"))
    
    resp, err := c.Client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("request failed: %v", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
    }
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response body: %v", err)
    }
    
    var result map[string]interface{}
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, fmt.Errorf("failed to parse JSON: %v", err)
    }
    
    return result, nil
}
// package utils

// import (
//     "encoding/json"
//     "fmt"
//     "io/ioutil"
//     "net/http"
//     "time"
//     "backend_rental/models"
// )

// func (c *ApiClient) FetchPropertyImages(hotelId string) (*models.PropertyImageResponse, error) {
//     url := fmt.Sprintf(
//         "https://booking-com18.p.rapidapi.com/stays/get-photos?hotelId=%s", 
//         hotelId,
//     )
    
//     req, _ := http.NewRequest("GET", url, nil)

//     for key, value := range c.Headers {
//         req.Header.Set(key, value)
//     }

//     client := &http.Client{
//         Timeout: 30 * time.Second,
//     }
//     resp, err := client.Do(req)
//     if err != nil {
//         return nil, err
//     }
//     defer resp.Body.Close()

//     body, _ := ioutil.ReadAll(resp.Body)
//     var imageResponse models.PropertyImageResponse
//     err = json.Unmarshal(body, &imageResponse)
//     if err != nil {
//         return nil, err
//     }

//     return &imageResponse, nil
// }