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

type PropertyImageService struct {
    RateLimiter *rate.Limiter
    ImageApiClient   *utils.ImageApiClient
    StoragePath string
    PropertiesPath string
}

func NewPropertyImageService() *PropertyImageService {
    limiter := rate.NewLimiter(rate.Every(6*time.Second), 10)
    
    dataDir := "data"
    if err := os.MkdirAll(dataDir, 0755); err != nil {
        fmt.Printf("Error creating data directory: %v\n", err)
    }
    
    return &PropertyImageService{
        RateLimiter: limiter,
        ImageApiClient:   utils.NewImageApiClient(),
        StoragePath: filepath.Join(dataDir, "property_images.json"),
        PropertiesPath: filepath.Join(dataDir, "properties.json"),
    }
}

func (s *PropertyImageService) LoadProperties() ([]models.Property, error) {
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

// func (s *PropertyImageService) FetchAllPropertyImages() ([]models.PropertyImage, error) {
//     properties, err := s.LoadProperties()
//     if err != nil {
//         return nil, err
//     }

//     var allPropertyImages []models.PropertyImage
//     ctx := context.Background()

//     for _, property := range properties {
//         startWait := time.Now()
//         err := s.RateLimiter.Wait(ctx)
//         waitDuration := time.Since(startWait)
        
//         fmt.Printf("Rate limiter wait for %s: %v\n", property.PropertyName, waitDuration)
        
//         if err != nil {
//             return nil, fmt.Errorf("rate limiter error: %v", err)
//         }

//         imageData, err := s.ApiClient.FetchPropertyImages(property.HotelID)
//         if err != nil {
//             fmt.Printf("Error fetching images for %s: %v\n", property.PropertyName, err)
//             continue
//         }

//         processedImages := s.processImages(property, imageData)
//         allPropertyImages = append(allPropertyImages, processedImages...)
//     }

//     fmt.Printf("Total property images fetched: %d\n", len(allPropertyImages))

//     // Save to file
//     err = s.SavePropertyImagesToFile(allPropertyImages)
//     if err != nil {
//         fmt.Printf("Warning: Failed to save property images: %v\n", err)
//     }

//     return allPropertyImages, nil
// }
func (s *PropertyImageService) FetchAllPropertyImages() ([]models.PropertyImage, error) {
    properties, err := s.LoadProperties()
    if err != nil {
        return nil, err
    }

    // Limit to first 10 properties
    if len(properties) > 10 {
        properties = properties[:10]
    }

    var allPropertyImages []models.PropertyImage
    ctx := context.Background()

    for _, property := range properties {
        startWait := time.Now()
        err := s.RateLimiter.Wait(ctx)
        waitDuration := time.Since(startWait)
        
        fmt.Printf("Fetching images for %s\n", property.PropertyName)
        
        if err != nil {
            return nil, fmt.Errorf("rate limiter error: %v", err)
        }

        imageData, err := s.ApiClient.FetchPropertyImages(property.HotelID)
        if err != nil {
            fmt.Printf("Error fetching images for %s: %v\n", property.PropertyName, err)
            continue
        }

        processedImages := s.processImages(property, imageData)
        allPropertyImages = append(allPropertyImages, processedImages...)
    }

    fmt.Printf("Total property images fetched: %d\n", len(allPropertyImages))

    // Save to file
    err = s.SavePropertyImagesToFile(allPropertyImages)
    if err != nil {
        fmt.Printf("Warning: Failed to save property images: %v\n", err)
    }

    return allPropertyImages, nil
}
func (s *PropertyImageService) processImages(property models.Property, imageData map[string]interface{}) []models.PropertyImage {
    var propertyImages []models.PropertyImage

    if data, ok := imageData["data"].(map[string]interface{}); ok {
        for hotelIDStr, hotelData := range data {
            hotelDataMap, ok := hotelData.(map[string]interface{})
            if !ok {
                continue
            }

            // Ensure hotel ID matches property
            hotelID, err := strconv.Atoi(hotelIDStr)
            if err != nil || hotelID != property.HotelID {
                continue
            }

            // Extract images for different tags
            tags := []string{"Property building", "Property", "Room"}
            for _, tag := range tags {
                imageUrls := s.extractImageUrlsForTag(hotelDataMap, tag)
                if len(imageUrls) > 0 {
                    propertyImages = append(propertyImages, models.PropertyImage{
                        PropertyID:   property.HotelID,
                        PropertyName: property.PropertyName,
                        Type:         tag,
                        Images:       imageUrls,
                    })
                }
            }
        }
    }

    return propertyImages
}

func (s *PropertyImageService) extractImageUrlsForTag(hotelData map[string]interface{}, targetTag string) []string {
    var imageUrls []string

    for _, entry := range hotelData {
        entryMap, ok := entry.(map[string]interface{})
        if !ok {
            continue
        }

        // Check tag
        tag, tagOk := entryMap["tag"].(string)
        if !tagOk || tag != targetTag {
            continue
        }

        // Extract images
        if images, ok := entryMap["images"].([]interface{}); ok {
            for _, img := range images {
                if url, ok := img.(string); ok {
                    imageUrls = append(imageUrls, url)
                }
            }
        }
    }

    return imageUrls
}
func (s *PropertyImageService) SavePropertyImagesToFile(images []models.PropertyImage) error {
    data, err := json.MarshalIndent(images, "", "    ")
    if err != nil {
        return fmt.Errorf("error marshaling property images data: %v", err)
    }

    err = os.WriteFile(s.StoragePath, data, 0644)
    if err != nil {
        return fmt.Errorf("error writing property images to file: %v", err)
    }

    fmt.Printf("Successfully saved %d property images to %s\n", len(images), s.StoragePath)
    return nil
}
// package services

// import (
//     "context"
//     "fmt"
//     "os"
//     "path/filepath"
//     "time"
//     "golang.org/x/time/rate"
//     "backend_rental/models"
//     "backend_rental/utils"
//     "encoding/json"
// 	"reflect"
// )

// type PropertyImageService struct {
//     RateLimiter *rate.Limiter
//     ApiClient   *utils.ApiClient
//     StoragePath string
//     PropertiesPath string
// }

// func NewPropertyImageService() *PropertyImageService {
//     limiter := rate.NewLimiter(rate.Every(6*time.Second), 10)
    
//     dataDir := "data"
//     if err := os.MkdirAll(dataDir, 0755); err != nil {
//         fmt.Printf("Error creating data directory: %v\n", err)
//     }
    
//     return &PropertyImageService{
//         RateLimiter: limiter,
//         ApiClient:   utils.NewApiClient(),
//         StoragePath: filepath.Join(dataDir, "property_images.json"),
//         PropertiesPath: filepath.Join(dataDir, "properties.json"),
//     }
// }

// func (s *PropertyImageService) LoadProperties() ([]models.Property, error) {
//     data, err := os.ReadFile(s.PropertiesPath)
//     if err != nil {
//         return nil, fmt.Errorf("error reading properties file: %v", err)
//     }

//     var properties []models.Property
//     err = json.Unmarshal(data, &properties)
//     if err != nil {
//         return nil, fmt.Errorf("error unmarshaling properties data: %v", err)
//     }

//     return properties, nil
// }



// func (s *PropertyImageService) LimitProperties(properties []models.Property, limit int) []models.Property {
//     if len(properties) <= limit {
//         return properties
//     }
//     return properties[:limit]
// }
// func (s *PropertyImageService) FetchPropertyImagesForProperties() ([]models.PropertyImage, error) {
//     properties, err := s.LoadProperties()
//     if err != nil {
//         return nil, fmt.Errorf("failed to load properties: %v", err)
//     }

//     // Limit to 10 properties
//     properties = s.LimitProperties(properties, 10)

//     var allPropertyImages []models.PropertyImage
//     ctx := context.Background()

//     for _, property := range properties {
//         startWait := time.Now()
//         err := s.RateLimiter.Wait(ctx)
//         waitDuration := time.Since(startWait)
        
//         fmt.Printf("Rate limiter wait for %s: %v\n", property.PropertyName, waitDuration)
        
//         if err != nil {
//             return nil, fmt.Errorf("rate limiter error: %v", err)
//         }

//         response, err := s.ApiClient.FetchPropertyImages(fmt.Sprintf("%d", property.HotelID))
//         if err != nil {
//             fmt.Printf("Error fetching images for %s: %v\n", property.PropertyName, err)
//             continue
//         }

//         // Enhanced debug logging
//         fmt.Printf("API Response for %s:\n", property.PropertyName)
//         fmt.Printf("Response Data Type: %T\n", response.Data)
//         fmt.Printf("Response Data Keys: %v\n", reflect.ValueOf(response.Data).MapKeys())
        
//         // Process images by type
//         propertyImages := s.processPropertyImages(property, response)
        
//         if len(propertyImages) == 0 {
//             fmt.Printf("No images found for property %s (Hotel ID: %d)\n", 
//                 property.PropertyName, property.HotelID)
//         }
        
//         allPropertyImages = append(allPropertyImages, propertyImages...)
//     }

//     if len(allPropertyImages) == 0 {
//         fmt.Println("Warning: No property images were found")
//     }

//     err = s.SavePropertyImagesToFile(allPropertyImages)
//     if err != nil {
//         return nil, fmt.Errorf("error saving property images: %v", err)
//     }

//     return allPropertyImages, nil
// }

// // func (s *PropertyImageService) FetchPropertyImagesForProperties() ([]models.PropertyImage, error) {
// //     properties, err := s.LoadProperties()
// //     if err != nil {
// //         return nil, err
// //     }

// //     var allPropertyImages []models.PropertyImage
// //     ctx := context.Background()

// //     for _, property := range properties {
// //         startWait := time.Now()
// //         err := s.RateLimiter.Wait(ctx)
// //         waitDuration := time.Since(startWait)
        
// //         fmt.Printf("Rate limiter wait for %s: %v\n", property.PropertyName, waitDuration)
        
// //         if err != nil {
// //             return nil, fmt.Errorf("rate limiter error: %v", err)
// //         }

// //         response, err := s.ApiClient.FetchPropertyImages(fmt.Sprintf("%d", property.HotelID))
// //         if err != nil {
// //             fmt.Printf("Error fetching images for %s: %v\n", property.PropertyName, err)
// //             continue
// //         }

// //         // Process images by type
// //         propertyImages := s.processPropertyImages(property, response)
// //         allPropertyImages = append(allPropertyImages, propertyImages...)
// //     }

// //     // Save to file
// //     err = s.SavePropertyImagesToFile(allPropertyImages)
// //     if err != nil {
// //         fmt.Printf("Warning: Failed to save property images: %v\n", err)
// //     }

// //     return allPropertyImages, nil
// // }


// func (s *PropertyImageService) processPropertyImages(property models.Property, response *models.PropertyImageResponse) []models.PropertyImage {
//     var propertyImages []models.PropertyImage

//     // Debug: Print entire response data
//     fmt.Printf("Full response data for property %s: %+v\n", property.PropertyName, response.Data)

//     // Convert to map with string keys
//     dataMap, ok := response.Data[fmt.Sprintf("%d", property.HotelID)]
//     if !ok {
//         fmt.Printf("No image data found for hotel ID %d\n", property.HotelID)
//         return propertyImages
//     }

//     // Debug: Print the specific data for this hotel
//     fmt.Printf("Hotel data: %+v\n", dataMap)

//     // Check the actual type of dataMap
//     switch v := dataMap.(type) {
//     case map[string]interface{}:
//         // Process as a map
//         imageTypes := []string{"Property building", "Property"}
//         for _, imageType := range imageTypes {
//             urls := s.extractImagesForType(map[string]interface{}{"data": v}, imageType)
//             if len(urls) > 0 {
//                 propertyImages = append(propertyImages, models.PropertyImage{
//                     PropertyID:   property.HotelID,
//                     PropertyName: property.PropertyName,
//                     ImageType:    imageType,
//                     ImageURLs:    urls,
//                 })
//             }
//         }
//     case []interface{}:
//         fmt.Printf("Unexpected slice type for hotel %d\n", property.HotelID)
//     default:
//         fmt.Printf("Unexpected type %T for hotel %d\n", dataMap, property.HotelID)
//     }

//     return propertyImages
// }

// func (s *PropertyImageService) extractImagesForType(imageData map[string]interface{}, imageType string) []string {
//     var urls []string
    
//     // More robust type checking and extraction
//     dataSlice, ok := imageData["data"].([]interface{})
//     if !ok {
//         fmt.Printf("Could not extract data slice for image type %s\n", imageType)
//         return urls
//     }

//     for _, item := range dataSlice {
//         itemMap, ok := item.(map[string]interface{})
//         if !ok {
//             continue
//         }

//         tag, tagOk := itemMap["tag"].(string)
//         if !tagOk || tag != imageType {
//             continue
//         }

//         // Extract image URLs
//         if imageList, ok := itemMap["4"].([]interface{}); ok {
//             for _, imgUrl := range imageList {
//                 if urlStr, ok := imgUrl.(string); ok {
//                     urls = append(urls, urlStr)
//                 }
//             }
//         }
//     }
    
//     return urls
// }


// func (s *PropertyImageService) extractImageURLs(imageData map[string]interface{}, imageType string) []string {
//     var urls []string
    
//     for _, item := range imageData {
//         itemMap, ok := item.(map[string]interface{})
//         if !ok {
//             continue
//         }
        
//         if tag, ok := itemMap["tag"].(string); ok && tag == imageType {
//             // Use "4" as a string key
//             if imageURLs, ok := itemMap["4"].([]interface{}); ok {
//                 for _, url := range imageURLs {
//                     if urlStr, ok := url.(string); ok {
//                         urls = append(urls, urlStr)
//                     }
//                 }
//             }
//         }
//     }
//     return urls
// }

// func (s *PropertyImageService) SavePropertyImagesToFile(propertyImages []models.PropertyImage) error {
//     data, err := json.MarshalIndent(propertyImages, "", "    ")
//     if err != nil {
//         return fmt.Errorf("error marshaling property images data: %v", err)
//     }

//     err = os.WriteFile(s.StoragePath, data, 0644)
//     if err != nil {
//         return fmt.Errorf("error writing property images to file: %v", err)
//     }

//     fmt.Printf("Successfully saved %d property images to %s\n", len(propertyImages), s.StoragePath)
//     return nil
// }