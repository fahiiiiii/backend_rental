package utils

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "time"
    // "backend_rental/models"
)

func (c *ApiClient) FetchPropertyDetails(hotelID int, checkIn, checkOut string) (*map[string]interface{}, error) {
    url := fmt.Sprintf(
        "https://booking-com18.p.rapidapi.com/stays/detail?hotelId=%d&checkinDate=%s&checkoutDate=%s&units=metric", 
        hotelID, 
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
    var propertyDetails map[string]interface{}
    err = json.Unmarshal(body, &propertyDetails)
    if err != nil {
        return nil, err
    }

    return &propertyDetails, nil
}
// package utils

// import (
//     "encoding/json"
//     "fmt"
//     "io/ioutil"
//     "net/http"
//     "time"
// )

// type PropertyDetailResponse struct {
//     HotelName               string `json:"hotel_name"`
//     AccommodationTypeName   string `json:"accommodation_type_name"`
//     BlockCount              int    `json:"block_count"`
//     PrivateBathroomCount    int    `json:"private_bathroom_count"`
//     Facilities              []struct {
//         Name string `json:"name"`
//     } `json:"facilities"`
// }

// func (c *ApiClient) FetchPropertyDetail(hotelId int, checkIn, checkOut string) (*PropertyDetailResponse, error) {
//     url := fmt.Sprintf(
//         "https://booking-com18.p.rapidapi.com/stays/detail?hotelId=%d&checkinDate=%s&checkoutDate=%s&units=metric", 
//         hotelId, 
//         checkIn, 
//         checkOut,
//     )
    
//     req, _ := http.NewRequest("GET", url, nil)

//     // Set headers
//     for key, value := range c.Headers {
//         req.Header.Set(key, value)
//     }

//     // Send request
//     client := &http.Client{
//         Timeout: 30 * time.Second,
//     }
//     resp, err := client.Do(req)
//     if err != nil {
//         return nil, err
//     }
//     defer resp.Body.Close()

//     // Parse response
//     body, _ := ioutil.ReadAll(resp.Body)
    
//     // Debug print the raw response
//     fmt.Println("Raw Response:", string(body))

//     var propertyDetailResponse PropertyDetailResponse
//     err = json.Unmarshal(body, &propertyDetailResponse)
//     if err != nil {
//         fmt.Println("Unmarshal Error:", err)
//         return nil, err
//     }

//     return &propertyDetailResponse, nil
// }
// --------------------------------------------------------
// package utils

// import (
//     "encoding/json"
//     "fmt"
//     "io/ioutil"
//     "net/http"
//     "time"
// )

// type PropertyDetailResponse struct {
//     HotelID                 int    `json:"hotel_id"`
//     HotelName               string `json:"hotel_name"`
//     AccommodationTypeName   string `json:"accommodation_type_name"`
//     BlockCount              int    `json:"block_count"`
//     PrivateBathroomCount    int    `json:"private_bathroom_count"`
//     Facilities              []struct {
//         Name string `json:"name"`
//     } `json:"facilities"`
// }

// func (c *ApiClient) FetchPropertyDetail(hotelId int, checkIn, checkOut string) (*PropertyDetailResponse, error) {
//     url := fmt.Sprintf(
//         "https://booking-com18.p.rapidapi.com/stays/detail?hotelId=%d&checkinDate=%s&checkoutDate=%s&units=metric", 
//         hotelId, 
//         checkIn, 
//         checkOut,
//     )
    
//     req, _ := http.NewRequest("GET", url, nil)

//     // Set headers
//     for key, value := range c.Headers {
//         req.Header.Set(key, value)
//     }

//     // Send request
//     client := &http.Client{
//         Timeout: 30 * time.Second,
//     }
//     resp, err := client.Do(req)
//     if err != nil {
//         return nil, err
//     }
//     defer resp.Body.Close()

//     // Parse response
//     body, _ := ioutil.ReadAll(resp.Body)
//     var propertyDetailResponse PropertyDetailResponse
//     err = json.Unmarshal(body, &propertyDetailResponse)
//     if err != nil {
//         return nil, err
//     }

//     return &propertyDetailResponse, nil
// }