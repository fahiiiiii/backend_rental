package services

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "time"
    "golang.org/x/time/rate"
    "backend_rental/models"
    "backend_rental/utils"
    "github.com/beego/beego/v2/client/orm"
)

type CityService struct {
    RateLimiter *rate.Limiter
    ApiClient   *utils.ApiClient
    StoragePath string
}

func NewCityService() *CityService {
    limiter := rate.NewLimiter(rate.Every(6*time.Second), 10)
    
    // Create a data directory if it doesn't exist
    dataDir := "data"
    if err := os.MkdirAll(dataDir, 0755); err != nil {
        fmt.Printf("Error creating data directory: %v\n", err)
    }
    
    return &CityService{
        RateLimiter: limiter,
        ApiClient:   utils.NewApiClient(),
        StoragePath: filepath.Join(dataDir, "cities.json"),
    }
}

// Save cities to both JSON file and database
func (s *CityService) SaveCities(cities []models.Location) error {
    // Save to JSON file
    if err := s.SaveCitiesToFile(cities); err != nil {
        return fmt.Errorf("failed to save to JSON: %v", err)
    }

    // Save to database
    if err := s.SaveCitiesToDB(cities); err != nil {
        return fmt.Errorf("failed to save to database: %v", err)
    }

    return nil
}

func (s *CityService) SaveCitiesToFile(cities []models.Location) error {
    data, err := json.MarshalIndent(cities, "", "    ")
    if err != nil {
        return fmt.Errorf("error marshaling cities data: %v", err)
    }

    err = os.WriteFile(s.StoragePath, data, 0644)
    if err != nil {
        return fmt.Errorf("error writing cities to file: %v", err)
    }

    fmt.Printf("Successfully saved %d cities to %s\n", len(cities), s.StoragePath)
    return nil
}

func (s *CityService) LoadCitiesFromFile() ([]models.Location, error) {
    data, err := os.ReadFile(s.StoragePath)
    if err != nil {
        if os.IsNotExist(err) {
            return nil, nil // Return empty slice if file doesn't exist
        }
        return nil, fmt.Errorf("error reading cities file: %v", err)
    }

    var cities []models.Location
    err = json.Unmarshal(data, &cities)
    if err != nil {
        return nil, fmt.Errorf("error unmarshaling cities data: %v", err)
    }

    return cities, nil
}

func (s *CityService) SaveCitiesToDB(cities []models.Location) error {
    o := orm.NewOrm()

    txOrm, err := o.Begin()
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %v", err)
    }

    defer func() {
        if r := recover(); r != nil {
            txOrm.Rollback()
            fmt.Printf("Transaction rolled back due to panic: %v\n", r)
        }
    }()

    query := "INSERT INTO location (city_name, city_id, country) VALUES "
    var values []string
    var params []interface{}
    
    for i, city := range cities {
        values = append(values, fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3))
        params = append(params, city.CityName, city.CityID, city.Country)
    }
    
    query += strings.Join(values, ",")
    _, err = txOrm.Raw(query, params...).Exec()
    if err != nil {
        txOrm.Rollback()
        return fmt.Errorf("failed to execute batch insert: %v", err)
    }

    err = txOrm.Commit()
    if err != nil {
        return fmt.Errorf("failed to commit transaction: %v", err)
    }

    fmt.Printf("Successfully saved %d cities to the database\n", len(cities))
    return nil
}

func (s *CityService) LoadCitiesFromDB() ([]models.Location, error) {
    o := orm.NewOrm()
    var cities []models.Location
    
    _, err := o.QueryTable("location").All(&cities)
    if err != nil {
        return nil, fmt.Errorf("error fetching cities from database: %v", err)
    }

    return cities, nil
}

func (s *CityService) FetchCitiesAlphabetically() ([]models.Location, error) {
    // Try loading from both sources
    dbCities, dbErr := s.LoadCitiesFromDB()
    fileCities, fileErr := s.LoadCitiesFromFile()

    // If database load was successful
    if dbErr == nil && len(dbCities) > 0 {
        fmt.Printf("Loaded %d cities from database\n", len(dbCities))
        return dbCities, nil
    }
    
    // If database failed but file load was successful
    if fileErr == nil && len(fileCities) > 0 {
        fmt.Printf("Database load failed (%v), but loaded %d cities from file\n", dbErr, len(fileCities))
        return fileCities, nil
    }

    // If both storage methods failed, log errors and fetch from API
    if dbErr != nil || fileErr != nil {
        fmt.Printf("Cache load failed - DB: %v, File: %v. Fetching from API...\n", dbErr, fileErr)
    }

    // Fetch from API
    var allCities []models.Location
    ctx := context.Background()

    for letter := 'A'; letter <= 'Z'; letter++ {
        query := string(letter)
        fmt.Printf("Fetching cities for query: %s\n", query)

        err := s.RateLimiter.Wait(ctx)
        if err != nil {
            return nil, fmt.Errorf("rate limiter error: %v", err)
        }

        response, err := s.ApiClient.FetchCityData(query)
        if err != nil {
            return nil, err
        }

        for _, item := range response.Data {
            if item.CityName != "" && item.CityID != "" && item.Country != "" {
                allCities = append(allCities, models.Location{
                    CityName: item.CityName,
                    CityID:   item.CityID,
                    Country:  item.Country,
                })
            }
        }
    }

    // Save the fetched data to both storage methods
    err := s.SaveCities(allCities)
    if err != nil {
        fmt.Printf("Warning: Failed to cache cities: %v\n", err)
    }

    return allCities, nil
}
