package services

import (
	// "context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	//"strings"
	// "sync"
	"time"
    "context"
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

// NewCityService initializes a new CityService with rate limiting, API client, and a storage path
func NewCityService() *CityService {
	limiter := rate.NewLimiter(rate.Every(6*time.Second), 10)

	dataDir := "data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		fmt.Printf("Error creating data directory: %v\n", err)
	} else {
		fmt.Printf("Data directory created/exists at: %s\n", dataDir)
	}

	absPath, err := filepath.Abs(dataDir)
	if err != nil {
		fmt.Printf("Error getting absolute path: %v\n", err)
	} else {
		fmt.Printf("Absolute path to data directory: %s\n", absPath)
	}

	return &CityService{
		RateLimiter: limiter,
		ApiClient:   utils.NewApiClient(),
		StoragePath: filepath.Join(dataDir, "cities.json"),
	}
}

// SaveCities saves cities to both a JSON file and a database
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
    if len(cities) == 0 {
        fmt.Println("No cities to save.")
        return nil
    }

    data, err := json.MarshalIndent(cities, "", "    ")
    if err != nil {
        return fmt.Errorf("error marshaling cities data: %v", err)
    }

    fmt.Printf("Writing %d cities to file: %s\n", len(cities), s.StoragePath)
    err = os.WriteFile(s.StoragePath, data, 0644)
    if err != nil {
        return fmt.Errorf("error writing cities to file: %v", err)
    }

    fmt.Printf("Successfully saved %d cities to %s\n", len(cities), s.StoragePath)
    return nil
}

// LoadCitiesFromFile loads city data from a JSON file
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
    if len(cities) == 0 {
        return nil
    }

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

    // Process cities in batches to avoid memory issues
    batchSize := 100
    for i := 0; i < len(cities); i += batchSize {
        end := i + batchSize
        if end > len(cities) {
            end = len(cities)
        }
        
        batch := cities[i:end]
        for _, city := range batch {
            // Check if city already exists
            existing := models.Location{CityID: city.CityID}
            err := txOrm.Read(&existing, "CityID")
            if err == orm.ErrNoRows {
                // City doesn't exist, insert it
                _, err := txOrm.Insert(&city)
                if err != nil {
                    txOrm.Rollback()
                    return fmt.Errorf("failed to insert city: %v", err)
                }
            }
        }
    }

    if err := txOrm.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %v", err)
    }

    fmt.Printf("Successfully saved %d cities to the database\n", len(cities))
    return nil
}


func (s *CityService) LoadCitiesFromDB() ([]models.Location, error) {
    o := orm.NewOrm()
    var cities []models.Location
    
    // Remove the default limit by specifying a large limit
    _, err := o.QueryTable("location").Limit(-1).All(&cities)
    if err != nil {
        return nil, fmt.Errorf("error fetching cities from database: %v", err)
    }

    fmt.Printf("Successfully loaded %d cities from database\n", len(cities))
    return cities, nil
}

func (s *CityService) FetchCitiesAlphabetically() ([]models.Location, error) {
    // Clear existing data first
    if err := s.clearExistingCities(); err != nil {
        fmt.Printf("Warning: Failed to clear existing cities: %v\n", err)
    }

    ctx := context.Background()
    rateLimiter := utils.NewRateLimiter(
        utils.LenientRateLimiter.Limit, 
        utils.LenientRateLimiter.BurstSize,
    )

    var allCities []models.Location
    
    // Sequential fetching with careful delays and detailed logging
    for letter := 'A'; letter <= 'Z'; letter++ {
        query := string(letter)
        fmt.Printf("\n=== Processing letter %s ===\n", query)

        // Wait for rate limit
        if err := utils.WaitForRateLimit(ctx, rateLimiter); err != nil {
            return nil, fmt.Errorf("rate limit error for letter %s: %v", query, err)
        }

        // Try up to 3 times with exponential backoff
        var response *models.ApiResponse
        var err error
        for attempt := 0; attempt < 3; attempt++ {
            fmt.Printf("Attempt %d: Fetching cities for letter %s...\n", attempt+1, query)
            
            sleepDuration := time.Second * time.Duration(3+attempt*2)
            time.Sleep(sleepDuration)

            response, err = s.ApiClient.FetchCityData(query)
            if err == nil {
                break
            }
            
            fmt.Printf("Attempt %d failed for letter %s: %v\n", attempt+1, query, err)
            if attempt < 2 {
                fmt.Printf("Waiting before retry...\n")
            }
        }

        if err != nil {
            fmt.Printf("Failed to fetch cities for letter %s after all attempts: %v\n", query, err)
            continue
        }

        if response == nil {
            fmt.Printf("No response received for letter %s\n", query)
            continue
        }

        // Log API response details
        fmt.Printf("API Response for letter %s: Total items received: %d\n", query, len(response.Data))

        // Process and save cities for this letter
        var letterCities []models.Location
        validCount := 0
        invalidCount := 0
        
        for _, item := range response.Data {
            if item.CityName != "" && item.CityID != "" && item.Country != "" {
                city := models.Location{
                    CityName: item.CityName,
                    CityID:   item.CityID,
                    Country:  item.Country,
                }
                letterCities = append(letterCities, city)
                allCities = append(allCities, city)
                validCount++
            } else {
                invalidCount++
                fmt.Printf("Skipped invalid city - Name: '%s', ID: '%s', Country: '%s'\n", 
                    item.CityName, item.CityID, item.Country)
            }
        }

        // Save progress after each successful letter
        if len(letterCities) > 0 {
            fmt.Printf("Letter %s summary:\n", query)
            fmt.Printf("- Valid cities found: %d\n", validCount)
            fmt.Printf("- Invalid entries skipped: %d\n", invalidCount)
            fmt.Printf("- Running total of all cities: %d\n", len(allCities))
            
            // Save current progress
            if err := s.SaveCities(allCities); err != nil {
                fmt.Printf("Warning: Failed to save progress: %v\n", err)
            } else {
                fmt.Printf("Successfully saved progress for letter %s\n", query)
            }
        } else {
            fmt.Printf("No valid cities found for letter %s\n", query)
        }

        sleepDuration := time.Second * 5
        fmt.Printf("Waiting %v before next letter...\n", sleepDuration)
        time.Sleep(sleepDuration)
    }

    if len(allCities) == 0 {
        return nil, fmt.Errorf("no cities were fetched from the API")
    }

    fmt.Printf("\n=== Final Summary ===\n")
    fmt.Printf("Successfully fetched a total of %d cities\n", len(allCities))
    return allCities, nil
}
func (s *CityService) clearExistingCities() error {
    o := orm.NewOrm()
    _, err := o.Raw("TRUNCATE TABLE location").Exec()
    if err != nil {
        return fmt.Errorf("failed to clear existing cities: %v", err)
    }
    
    // Also remove the JSON file if it exists
    if err := os.Remove(s.StoragePath); err != nil && !os.IsNotExist(err) {
        return fmt.Errorf("failed to remove cities file: %v", err)
    }
    
    fmt.Println("Successfully cleared existing cities from database and file")
    return nil
}

