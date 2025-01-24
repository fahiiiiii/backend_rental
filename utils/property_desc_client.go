package utils

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "time"
)

type PropertyDescriptionResponse struct {
    Data []PropertyDescription `json:"data"`
    Status interface{} `json:"status"`
    Message interface{} `json:"message"`
}

type PropertyDescription struct {
    Description     string `json:"description"`
    DescriptionTypeID int `json:"descriptiontype_id"`
    LanguageCode    string `json:"languagecode"`
}

func (c *ApiClient) FetchPropertyDescription(hotelID string) (*PropertyDescriptionResponse, error) {
    url := fmt.Sprintf(
        "https://booking-com18.p.rapidapi.com/stays/get-description?hotelId=%s", 
        hotelID,
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
    var descResponse PropertyDescriptionResponse
    err = json.Unmarshal(body, &descResponse)
    if err != nil {
        // Print raw response for debugging
        fmt.Printf("Raw response: %s\n", string(body))
        return nil, fmt.Errorf("error unmarshaling response: %v", err)
    }

    return &descResponse, nil
}
// package utils

// import (
//     "encoding/json"
//     "fmt"
//     "io/ioutil"
//     "net/http"
//     "time"
//     // "backend_rental/models"
// )

// type PropertyDescriptionResponse struct {
//     Data []PropertyDescription `json:"data"`
//     Status string `json:"status"`
//     Message string `json:"message"`
// }

// type PropertyDescription struct {
//     Description     string `json:"description"`
//     DescriptionTypeID int `json:"descriptiontype_id"`
//     LanguageCode    string `json:"languagecode"`
// }

// func (c *ApiClient) FetchPropertyDescription(hotelID string) (*PropertyDescriptionResponse, error) {
//     url := fmt.Sprintf(
//         "https://booking-com18.p.rapidapi.com/stays/get-description?hotelId=%s", 
//         hotelID,
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
//     var descResponse PropertyDescriptionResponse
//     err = json.Unmarshal(body, &descResponse)
//     if err != nil {
//         return nil, err
//     }

//     return &descResponse, nil
// }