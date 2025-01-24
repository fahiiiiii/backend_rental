
package services

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "time"
    "golang.org/x/time/rate"
    "backend_rental/models"
    "backend_rental/utils"
)

type PropertyDetailsService struct {
    RateLimiter *rate.Limiter
    ApiClient   *utils.ApiClient
    StoragePath string
    PropertiesPath string
}

func NewPropertyDetailsService() *PropertyDetailsService {
    limiter := rate.NewLimiter(rate.Every(6*time.Second), 10)
    
    dataDir := "data"
    if err := os.MkdirAll(dataDir, 0755); err != nil {
        fmt.Printf("Error creating data directory: %v\n", err)
    }
    
    return &PropertyDetailsService{
        RateLimiter: limiter,
        ApiClient:   utils.NewApiClient(),
        StoragePath: filepath.Join(dataDir, "property_details.json"),
        PropertiesPath: filepath.Join(dataDir, "properties.json"),
    }
}

func (s *PropertyDetailsService) LoadProperties() ([]models.Property, error) {
    data, err := os.ReadFile(s.PropertiesPath)
    if err != nil {
        return nil, fmt.Errorf("error reading properties file: %v", err)
    }

    var properties []models.Property
    err = json.Unmarshal(data, &properties)
    if err != nil {
        return nil, fmt.Errorf("error unmarshaling properties data: %v", err)
    }

    return properties, nil
}
func (s *PropertyDetailsService) FetchPropertyDetails() ([]models.PropertyDetail, error) {
    properties, err := s.LoadProperties()
    if err != nil {
        return nil, err
    }

    var allPropertyDetails []models.PropertyDetail
    ctx := context.Background()
    checkIn := time.Now().AddDate(0, 0, 30).Format("2006-01-02")
    checkOut := time.Now().AddDate(0, 0, 31).Format("2006-01-02")

    // Limit to 10 properties
    limit := 10
    if len(properties) < limit {
        limit = len(properties)
    }

    for _, property := range properties[:limit] {
        startWait := time.Now()
        err := s.RateLimiter.Wait(ctx)
        waitDuration := time.Since(startWait)
        
        fmt.Printf("Rate limiter wait for %s: %v\n", property.PropertyName, waitDuration)
        
        if err != nil {
            return nil, fmt.Errorf("rate limiter error: %v", err)
        }

        response, err := s.ApiClient.FetchPropertyDetails(property.HotelID, checkIn, checkOut)
        if err != nil {
            fmt.Printf("Error fetching details for %s (ID: %d): %v\n", property.PropertyName, property.HotelID, err)
            continue
        }

        // Debug print the entire response
        fmt.Printf("Full API Response for %s: %+v\n", property.PropertyName, response)

        // Extract required details with explicit type checking
        propertyDetail := models.PropertyDetail{
			HotelID:       property.HotelID,
			CityID:        property.CityID,
			PropertyType:  s.extractPropertyType(response),
			Bedrooms:      s.extractBedrooms(response),
			Bathrooms:     s.extractBathrooms(response),
			Amenities:     s.extractAmenities(response),
		}
        allPropertyDetails = append(allPropertyDetails, propertyDetail)
    }

    fmt.Printf("Property details fetched: %d\n", len(allPropertyDetails))

    // Save to file
    err = s.SavePropertyDetailsToFile(allPropertyDetails)
    if err != nil {
        fmt.Printf("Warning: Failed to save property details: %v\n", err)
    }

    return allPropertyDetails, nil
}
// func (s *PropertyDetailsService) FetchPropertyDetails() ([]models.PropertyDetail, error) {
//     properties, err := s.LoadProperties()
//     if err != nil {
//         return nil, err
//     }

//     var allPropertyDetails []models.PropertyDetail
//     ctx := context.Background()
//     checkIn := time.Now().AddDate(0, 0, 30).Format("2006-01-02")
//     checkOut := time.Now().AddDate(0, 0, 31).Format("2006-01-02")

//     for _, property := range properties {
//         startWait := time.Now()
//         err := s.RateLimiter.Wait(ctx)
//         waitDuration := time.Since(startWait)
        
//         fmt.Printf("Rate limiter wait for %s: %v\n", property.PropertyName, waitDuration)
        
//         if err != nil {
//             return nil, fmt.Errorf("rate limiter error: %v", err)
//         }

//         response, err := s.ApiClient.FetchPropertyDetails(property.HotelID, checkIn, checkOut)
//         if err != nil {
//             fmt.Printf("Error fetching details for %s (ID: %d): %v\n", property.PropertyName, property.HotelID, err)
//             continue
//         }

//         // Extract required details
//         propertyDetail := models.PropertyDetail{
//             HotelID:     property.HotelID,
//             CityID:      property.CityID,
//             Description: s.extractDescription(response),
//             PropertyType: s.extractPropertyType(response),
//             Bedrooms:    s.extractBedrooms(response),
//             Bathrooms:   s.extractBathrooms(response),
//             Amenities:   s.extractAmenities(response),
//             Address:     s.extractAddress(response),
//             HotelName:   s.extractHotelName(response),
//         }

//         allPropertyDetails = append(allPropertyDetails, propertyDetail)
//     }

//     fmt.Printf("Total property details fetched: %d\n", len(allPropertyDetails))

//     // Save to file
//     err = s.SavePropertyDetailsToFile(allPropertyDetails)
//     if err != nil {
//         fmt.Printf("Warning: Failed to save property details: %v\n", err)
//     }

//     return allPropertyDetails, nil
// }

func (s *PropertyDetailsService) extractPropertyType(details *map[string]interface{}) string {
    if data, ok := (*details)["data"].(map[string]interface{}); ok {
        if typeName, ok := data["accommodation_type_name"].(string); ok {
            return typeName
        }
    }
    fmt.Println("Failed to extract accommodation_type_name")
    return ""
}

func (s *PropertyDetailsService) extractBedrooms(details *map[string]interface{}) int {
    if data, ok := (*details)["data"].(map[string]interface{}); ok {
        if blocks, ok := data["block_count"].(float64); ok {
            return int(blocks)
        }
    }
    fmt.Println("Failed to extract block_count")
    return 0
}

func (s *PropertyDetailsService) extractBathrooms(details *map[string]interface{}) int {
    if data, ok := (*details)["data"].(map[string]interface{}); ok {
        // Try multiple potential ways to extract bathroom count
        if bathrooms, ok := data["number_of_bathrooms"].(float64); ok && bathrooms > 0 {
            return int(bathrooms)
        }
        // Fallback to block count if no specific bathroom count
        if blocks, ok := data["block_count"].(float64); ok {
            return int(blocks)
        }
    }
    return 0
}

func (s *PropertyDetailsService) extractAmenities(details *map[string]interface{}) []string {
    var amenities []string
    if data, ok := (*details)["data"].(map[string]interface{}); ok {
        // Try multiple sources for amenities
        if facilities, ok := data["facilities"].([]interface{}); ok {
            for _, facility := range facilities {
                if facilityMap, ok := facility.(map[string]interface{}); ok {
                    if name, ok := facilityMap["name"].(string); ok {
                        amenities = append(amenities, name)
                    }
                }
            }
        }
        
        // Fallback to facilities_block
        if facilitiesBlock, ok := data["facilities_block"].(map[string]interface{}); ok {
            if facilities, ok := facilitiesBlock["facilities"].([]interface{}); ok {
                for _, facility := range facilities {
                    if facilityMap, ok := facility.(map[string]interface{}); ok {
                        if name, ok := facilityMap["name"].(string); ok {
                            amenities = append(amenities, name)
                        }
                    }
                }
            }
        }
    }
    
    return amenities
}

func (s *PropertyDetailsService) SavePropertyDetailsToFile(propertyDetails []models.PropertyDetail) error {
    data, err := json.MarshalIndent(propertyDetails, "", "    ")
    if err != nil {
        return fmt.Errorf("error marshaling property details data: %v", err)
    }

    err = os.WriteFile(s.StoragePath, data, 0644)
    if err != nil {
        return fmt.Errorf("error writing property details to file: %v", err)
    }

    fmt.Printf("Successfully saved %d property details to %s\n", len(propertyDetails), s.StoragePath)
    return nil
}

// package services

// import (
//     "context"
//     "encoding/json"
//     "fmt"
//     "os"
//     "path/filepath"
//     "time"
//     "golang.org/x/time/rate"
//     "backend_rental/models"
//     "backend_rental/utils"
// )

// type PropertyDetailService struct {
//     RateLimiter *rate.Limiter
//     ApiClient   *utils.ApiClient
//     StoragePath string
// }

// func NewPropertyDetailService() *PropertyDetailService {
//     limiter := rate.NewLimiter(rate.Every(6*time.Second), 10)
    
//     dataDir := "data"
//     if err := os.MkdirAll(dataDir, 0755); err != nil {
//         fmt.Printf("Error creating data directory: %v\n", err)
//     }
    
//     return &PropertyDetailService{
//         RateLimiter: limiter,
//         ApiClient:   utils.NewApiClient(),
//         StoragePath: filepath.Join(dataDir, "property_details.json"),
//     }
// }

// func (s *PropertyDetailService) FetchPropertyDetails() ([]models.PropertyDetail, error) {
//     propertiesData, err := os.ReadFile("data/properties.json")
//     if err != nil {
//         return nil, fmt.Errorf("error reading properties file: %v", err)
//     }

//     var properties []models.Property
//     err = json.Unmarshal(propertiesData, &properties)
//     if err != nil {
//         return nil, fmt.Errorf("error unmarshaling properties data: %v", err)
//     }

//     var allPropertyDetails []models.PropertyDetail
//     ctx := context.Background()
//     checkIn := time.Now().AddDate(0, 0, 30).Format("2006-01-02")
//     checkOut := time.Now().AddDate(0, 0, 31).Format("2006-01-02")

//     // Limit to first 10 properties
//     limit := 10
//     if len(properties) < limit {
//         limit = len(properties)
//     }

//     for _, property := range properties[:limit] {
//         startWait := time.Now()
//         err := s.RateLimiter.Wait(ctx)
//         waitDuration := time.Since(startWait)
        
//         fmt.Printf("Rate limiter wait for hotel %d: %v\n", property.HotelID, waitDuration)
        
//         if err != nil {
//             return nil, fmt.Errorf("rate limiter error: %v", err)
//         }

//         response, err := s.ApiClient.FetchPropertyDetail(property.HotelID, checkIn, checkOut)
//         if err != nil {
//             fmt.Printf("Error fetching details for hotel %d: %v\n", property.HotelID, err)
//             continue
//         }

//         propertyDetail := models.PropertyDetail{
//             PropertyID:    property.HotelID,
//             HotelID:       property.HotelID,
//             CityID:        property.CityID,
//             Name:          response.HotelName,
//             PropertyType:  response.AccommodationTypeName,
//             Bedrooms:      response.BlockCount,
//             Bathrooms:     response.PrivateBathroomCount,
//             Amenities:     extractAmenities(response.Facilities),
//         }

//         allPropertyDetails = append(allPropertyDetails, propertyDetail)
//     }

//     fmt.Printf("Total property details fetched: %d\n", len(allPropertyDetails))

//     // Save to file
//     err = s.SavePropertyDetailsToFile(allPropertyDetails)
//     if err != nil {
//         fmt.Printf("Warning: Failed to save property details: %v\n", err)
//     }

//     return allPropertyDetails, nil
// }
// // func (s *PropertyDetailService) FetchPropertyDetails() ([]models.PropertyDetail, error) {
// //     // Read properties from file
// //     propertiesData, err := os.ReadFile("data/properties.json")
// //     if err != nil {
// //         return nil, fmt.Errorf("error reading properties file: %v", err)
// //     }

// //     var properties []models.Property
// //     err = json.Unmarshal(propertiesData, &properties)
// //     if err != nil {
// //         return nil, fmt.Errorf("error unmarshaling properties data: %v", err)
// //     }

// //     var allPropertyDetails []models.PropertyDetail
// //     ctx := context.Background()
// //     checkIn := time.Now().AddDate(0, 0, 30).Format("2006-01-02")
// //     checkOut := time.Now().AddDate(0, 0, 31).Format("2006-01-02")

// //     for _, property := range properties {
// //         startWait := time.Now()
// //         err := s.RateLimiter.Wait(ctx)
// //         waitDuration := time.Since(startWait)
        
// //         fmt.Printf("Rate limiter wait for hotel %d: %v\n", property.HotelID, waitDuration)
        
// //         if err != nil {
// //             return nil, fmt.Errorf("rate limiter error: %v", err)
// //         }

// //         response, err := s.ApiClient.FetchPropertyDetail(property.HotelID, checkIn, checkOut)
// //         if err != nil {
// //             fmt.Printf("Error fetching details for hotel %d: %v\n", property.HotelID, err)
// //             continue
// //         }

// //         // Verify the response contains data
// //         fmt.Printf("Hotel Name: %s, Accommodation Type: %s\n", 
// //             response.HotelName, response.AccommodationTypeName)

// //         propertyDetail := models.PropertyDetail{
// //             PropertyID:    property.HotelID,  // Use HotelID as PropertyID
// //             HotelID:       property.HotelID,
// //             CityID:        property.CityID,
// //             Name:          response.HotelName,
// //             PropertyType:  response.AccommodationTypeName,
// //             Bedrooms:      response.BlockCount,
// //             Bathrooms:     response.PrivateBathroomCount,
// //             Amenities:     extractAmenities(response.Facilities),
// //         }

// //         allPropertyDetails = append(allPropertyDetails, propertyDetail)
// //     }

// //     fmt.Printf("Total property details fetched: %d\n", len(allPropertyDetails))

// //     // Save to file
// //     err = s.SavePropertyDetailsToFile(allPropertyDetails)
// //     if err != nil {
// //         fmt.Printf("Warning: Failed to save property details: %v\n", err)
// //     }

// //     return allPropertyDetails, nil
// // }

// func extractAmenities(facilities []struct {
//     Name string `json:"name"`
// }) []string {
//     var amenities []string
//     for _, facility := range facilities {
//         if facility.Name != "" {
//             amenities = append(amenities, facility.Name)
//         }
//     }
//     return amenities
// }

// func (s *PropertyDetailService) SavePropertyDetailsToFile(propertyDetails []models.PropertyDetail) error {
//     data, err := json.MarshalIndent(propertyDetails, "", "    ")
//     if err != nil {
//         return fmt.Errorf("error marshaling property details data: %v", err)
//     }

//     err = os.WriteFile(s.StoragePath, data, 0644)
//     if err != nil {
//         return fmt.Errorf("error writing property details to file: %v", err)
//     }

//     fmt.Printf("Successfully saved %d property details to %s\n", len(propertyDetails), s.StoragePath)
//     return nil
// }

