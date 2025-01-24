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
    RateLimiter   *rate.Limiter
    ApiClient     *utils.ApiClient
    StoragePath   string
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
        ApiClient:   utils.NewApiClient(),
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

func (s *PropertyImageService) FetchImagesForProperties() ([]models.PropertyImage, error) {
    properties, err := s.LoadProperties()
    if err != nil {
        return nil, err
    }

    var allPropertyImages []models.PropertyImage
    ctx := context.Background()

    for _, property := range properties {
        startWait := time.Now()
        err := s.RateLimiter.Wait(ctx)
        waitDuration := time.Since(startWait)
        
        fmt.Printf("Rate limiter wait for %s: %v\n", property.PropertyName, waitDuration)
        
        if err != nil {
            return nil, fmt.Errorf("rate limiter error: %v", err)
        }

        response, err := s.ApiClient.FetchPropertyImages(fmt.Sprintf("%d", property.HotelID))
        if err != nil {
            fmt.Printf("Error fetching images for %s: %v\n", property.PropertyName, err)
            continue
        }

        // Process image types based on the response structure
        for typeName, data := range response.Data {
            imageUrls := extractImageUrls(data)
            if len(imageUrls) > 0 {
                propertyImage := models.PropertyImage{
                    PropertyID:   property.HotelID,
                    PropertyName: property.PropertyName,
                    Type:         typeName,
                    URLs:         imageUrls,
                }
                allPropertyImages = append(allPropertyImages, propertyImage)
            }
        }
    }

    fmt.Printf("Total property images fetched: %d\n", len(allPropertyImages))

    err = s.SavePropertyImagesToFile(allPropertyImages)
    if err != nil {
        fmt.Printf("Warning: Failed to save property images: %v\n", err)
    }

    return allPropertyImages, nil
}

func extractImageUrls(data interface{}) []string {
    // Implement logic to extract image URLs based on the API response structure
    // This will depend on the exact structure of the response
    // You might need to use type assertions and recursive parsing
    // Example placeholder implementation
    urls := []string{}
    switch v := data.(type) {
    case map[string]interface{}:
        for _, value := range v {
            if urlStr, ok := value.(string); ok && isImageUrl(urlStr) {
                urls = append(urls, urlStr)
            }
        }
    case []interface{}:
        for _, item := range v {
            if urlStr, ok := item.(string); ok && isImageUrl(urlStr) {
                urls = append(urls, urlStr)
            }
        }
    }
    return urls
}

func isImageUrl(url string) bool {
    // Basic URL validation for image URLs
    return len(url) > 0 && (strings.Contains(url, ".jpg") || 
                             strings.Contains(url, ".png") || 
                             strings.Contains(url, ".jpeg"))
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