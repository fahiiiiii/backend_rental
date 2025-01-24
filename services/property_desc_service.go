
package services

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "strconv"
    "time"
    "golang.org/x/time/rate"
    "backend_rental/models"
    "backend_rental/utils"
)

type PropertyDescService struct {
    RateLimiter *rate.Limiter
    ApiClient   *utils.ApiClient
    StoragePath string
}

type PropertyDescriptionDetail struct {
    PropertyID   int    `json:"property_id"`
    PropertyName string `json:"property_name"`
    Description  string `json:"description"`
}

func NewPropertyDescService() *PropertyDescService {
    limiter := rate.NewLimiter(rate.Every(6*time.Second), 10)
    
    dataDir := "data"
    if err := os.MkdirAll(dataDir, 0755); err != nil {
        fmt.Printf("Error creating data directory: %v\n", err)
    }
    
    return &PropertyDescService{
        RateLimiter: limiter,
        ApiClient:   utils.NewApiClient(),
        StoragePath: filepath.Join(dataDir, "property_desc_image.json"),
    }
}
// Modify the FetchAndSavePropertyDescriptions method
func (s *PropertyDescService) FetchAndSavePropertyDescriptions() error {
    // Load properties from file
    propertiesData, err := os.ReadFile("data/properties.json")
    if err != nil {
        return fmt.Errorf("error reading properties file: %v", err)
    }

    var properties []models.Property
    err = json.Unmarshal(propertiesData, &properties)
    if err != nil {
        return fmt.Errorf("error unmarshaling properties: %v", err)
    }

    // Limit to first 10 properties
    if len(properties) > 10 {
        properties = properties[:10]
    }

    var propertyDescriptions []PropertyDescriptionDetail
    ctx := context.Background()

    for _, property := range properties {
        // Rate limit
        startWait := time.Now()
        err := s.RateLimiter.Wait(ctx)
        waitDuration := time.Since(startWait)
        
        fmt.Printf("Rate limiter wait for %s: %v\n", property.PropertyName, waitDuration)
        
        if err != nil {
            return fmt.Errorf("rate limiter error: %v", err)
        }

        // Fetch description
        response, err := s.ApiClient.FetchPropertyDescription(strconv.Itoa(property.HotelID))
        if err != nil {
            fmt.Printf("Error fetching description for %s: %v\n", property.PropertyName, err)
            continue
        }

        // Find the main description (typically the first one with descriptiontype_id 6)
        var mainDescription string
        for _, desc := range response.Data {
            if desc.DescriptionTypeID == 6 {
                mainDescription = desc.Description
                break
            }
        }

        // Add to results
        propertyDescriptions = append(propertyDescriptions, PropertyDescriptionDetail{
            PropertyID:   property.HotelID,
            PropertyName: property.PropertyName,
            Description:  mainDescription,
        })

        fmt.Printf("Fetched description for %s\n", property.PropertyName)
    }

    // Save to file
    return s.SavePropertyDescriptionsToFile(propertyDescriptions)
}
// func (s *PropertyDescService) FetchAndSavePropertyDescriptions() error {
//     // Load properties from file
//     propertiesData, err := os.ReadFile("data/properties.json")
//     if err != nil {
//         return fmt.Errorf("error reading properties file: %v", err)
//     }

//     var properties []models.Property
//     err = json.Unmarshal(propertiesData, &properties)
//     if err != nil {
//         return fmt.Errorf("error unmarshaling properties: %v", err)
//     }

//     var propertyDescriptions []PropertyDescriptionDetail
//     ctx := context.Background()

//     for _, property := range properties {
//         // Rate limit
//         startWait := time.Now()
//         err := s.RateLimiter.Wait(ctx)
//         waitDuration := time.Since(startWait)
        
//         fmt.Printf("Rate limiter wait for %s: %v\n", property.PropertyName, waitDuration)
        
//         if err != nil {
//             return fmt.Errorf("rate limiter error: %v", err)
//         }

//         // Fetch description
//         response, err := s.ApiClient.FetchPropertyDescription(strconv.Itoa(property.HotelID))
//         if err != nil {
//             fmt.Printf("Error fetching description for %s: %v\n", property.PropertyName, err)
//             continue
//         }

//         // Find the main description (typically the first one with descriptiontype_id 6)
//         var mainDescription string
//         for _, desc := range response.Data {
//             if desc.DescriptionTypeID == 6 {
//                 mainDescription = desc.Description
//                 break
//             }
//         }

//         // Add to results
//         propertyDescriptions = append(propertyDescriptions, PropertyDescriptionDetail{
//             PropertyID:   property.HotelID,
//             PropertyName: property.PropertyName,
//             Description:  mainDescription,
//         })

//         fmt.Printf("Fetched description for %s\n", property.PropertyName)
//     }

//     // Save to file
//     return s.SavePropertyDescriptionsToFile(propertyDescriptions)
// }

func (s *PropertyDescService) SavePropertyDescriptionsToFile(descriptions []PropertyDescriptionDetail) error {
    data, err := json.MarshalIndent(descriptions, "", "    ")
    if err != nil {
        return fmt.Errorf("error marshaling descriptions data: %v", err)
    }

    err = os.WriteFile(s.StoragePath, data, 0644)
    if err != nil {
        return fmt.Errorf("error writing descriptions to file: %v", err)
    }

    fmt.Printf("Successfully saved %d property descriptions to %s\n", len(descriptions), s.StoragePath)
    return nil
}