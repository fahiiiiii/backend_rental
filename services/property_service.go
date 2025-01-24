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

type PropertyService struct {
    RateLimiter *rate.Limiter
    ApiClient   *utils.ApiClient
    StoragePath string
    CitiesPath  string
}

func NewPropertyService() *PropertyService {
    limiter := rate.NewLimiter(rate.Every(6*time.Second), 10)
    
    dataDir := "data"
    if err := os.MkdirAll(dataDir, 0755); err != nil {
        fmt.Printf("Error creating data directory: %v\n", err)
    }
    
    return &PropertyService{
        RateLimiter: limiter,
        ApiClient:   utils.NewApiClient(),
        StoragePath: filepath.Join(dataDir, "properties.json"),
        CitiesPath:  filepath.Join(dataDir, "cities.json"),
    }
}

func (s *PropertyService) LoadCities() ([]models.Location, error) {
    data, err := os.ReadFile(s.CitiesPath)
    if err != nil {
        return nil, fmt.Errorf("error reading cities file: %v", err)
    }

    var cities []models.Location
    err = json.Unmarshal(data, &cities)
    if err != nil {
        return nil, fmt.Errorf("error unmarshaling cities data: %v", err)
    }

    return cities, nil
}
func (s *PropertyService) FetchPropertiesForCities() ([]models.Property, error) {
    cities, err := s.LoadCities()
    if err != nil {
        return nil, err
    }

    var allProperties []models.Property
    ctx := context.Background()
    checkIn := time.Now().AddDate(0, 0, 30).Format("2006-01-02")
    checkOut := time.Now().AddDate(0, 0, 31).Format("2006-01-02")

    for _, city := range cities {
        startWait := time.Now()
        err := s.RateLimiter.Wait(ctx)
        waitDuration := time.Since(startWait)
        
        fmt.Printf("Rate limiter wait for %s: %v\n", city.CityName, waitDuration)
        
        if err != nil {
            return nil, fmt.Errorf("rate limiter error: %v", err)
        }

        response, err := s.ApiClient.FetchPropertiesForCity(city.CityID, checkIn, checkOut)
        if err != nil {
            fmt.Printf("Error fetching properties for %s: %v\n", city.CityName, err)
            continue
        }

        fmt.Printf("Fetched %d properties for %s\n", len(response.Data), city.CityName)

        for _, property := range response.Data {
            property.CityID = city.CityID
            allProperties = append(allProperties, property)
        }
    }

    fmt.Printf("Total properties fetched: %d\n", len(allProperties))

    // Save to file
    err = s.SavePropertiesToFile(allProperties)
    if err != nil {
        fmt.Printf("Warning: Failed to save properties: %v\n", err)
    }

    return allProperties, nil
}

func (s *PropertyService) SavePropertiesToFile(properties []models.Property) error {
    data, err := json.MarshalIndent(properties, "", "    ")
    if err != nil {
        return fmt.Errorf("error marshaling properties data: %v", err)
    }

    err = os.WriteFile(s.StoragePath, data, 0644)
    if err != nil {
        return fmt.Errorf("error writing properties to file: %v", err)
    }

    fmt.Printf("Successfully saved %d properties to %s\n", len(properties), s.StoragePath)
    return nil
}