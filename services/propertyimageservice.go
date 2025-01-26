package services

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "time"
    "strconv"
    "backend_rental/models"
    "backend_rental/utils"
    "golang.org/x/time/rate"
)

type PropertyImageService struct {
    apiBaseURL   string
    apiKey       string
    rateLimiter  *rate.Limiter
}

func NewPropertyImageService(apiBaseURL, apiKey string) (*PropertyImageService, error) {
    return &PropertyImageService{
        apiBaseURL:  apiBaseURL,
        apiKey:      apiKey,
        rateLimiter: utils.NewRateLimiter(
            utils.LenientRateLimiter.Limit, 
            utils.LenientRateLimiter.BurstSize,
        ),
    }, nil
}

func (s *PropertyImageService) FetchPropertyImages(ctx context.Context, properties []models.Property) ([]models.PropertyImage, error) {
    var allPropertyImages []models.PropertyImage
    
    // Limit to first 10 properties
    if len(properties) > 10 {
        properties = properties[:10]
    }

    log.Printf("Fetching images for %d properties", len(properties))

    for i, property := range properties {
        log.Printf("Fetching images for property %d (ID: %d)", i+1, property.HotelID)

        if err := s.rateLimiter.Wait(ctx); err != nil {
            log.Printf("Rate limiter error: %v", err)
            return nil, err
        }

        // Modify URL to use string ID format from curl example
        url := fmt.Sprintf("https://%s/web/stays/details?id=us/mayfair-new-york&checkIn=%s&checkOut=%s", 
            s.apiBaseURL, 
            time.Now().AddDate(0, 1, 0).Format("2006-01-02"),
            time.Now().AddDate(0, 1, 7).Format("2006-01-02"),
        )

        req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
        if err != nil {
            log.Printf("Error creating request for property %d: %v", property.HotelID, err)
            continue
        }
        req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
        req.Header.Add("x-rapidapi-key", s.apiKey)

        resp, err := http.DefaultClient.Do(req)
        if err != nil {
            log.Printf("HTTP request error for property %d: %v", property.HotelID, err)
            continue
        }
        defer resp.Body.Close()

        // Check response status code
        if resp.StatusCode != http.StatusOK {
            log.Printf("Unexpected status code for property %d: %d", property.HotelID, resp.StatusCode)
            body, _ := io.ReadAll(resp.Body)
            log.Printf("Response body: %s", string(body))
            continue
        }

        body, err := io.ReadAll(resp.Body)
        if err != nil {
            log.Printf("Error reading response body for property %d: %v", property.HotelID, err)
            continue
        }

        var apiResponse map[string]interface{}
        if err := json.Unmarshal(body, &apiResponse); err != nil {
            log.Printf("JSON unmarshaling error for property %d: %v", property.HotelID, err)
            log.Printf("Response body: %s", string(body))
            continue
        }

        // More detailed logging for debugging
        var keys []string
        for k := range apiResponse {
            keys = append(keys, k)
        }
        log.Printf("API Response keys for property %d: %v", property.HotelID, keys)

        data, ok := apiResponse["data"].(map[string]interface{})
        if !ok {
            log.Printf("No 'data' key found for property %d", property.HotelID)
            log.Printf("Full API response: %+v", apiResponse)
            continue
        }

        // Hotel Photos
        hotelPhotos, ok := data["hotelPhotos"].([]interface{})
        if ok {
            var hotelImageURLs []string
            for _, photo := range hotelPhotos {
                photoMap, ok := photo.(map[string]interface{})
                if !ok {
                    continue
                }
                thumbURL, ok := photoMap["thumb_url"].(string)
                if !ok {
                    continue
                }
                hotelImageURLs = append(hotelImageURLs, thumbURL)
            }

            if len(hotelImageURLs) > 0 {
                hotelPropertyImage := models.PropertyImage{
                    PropertyID:   property.HotelID,
                    PropertyName: strconv.Itoa(property.HotelID),
                    ImageType:    "hotel_photos",
                    ImageURLs:    hotelImageURLs,
                }
                allPropertyImages = append(allPropertyImages, hotelPropertyImage)
            }
        }

        // Room Photos
        roomPhotos, ok := data["allRoomPhotos"].([]interface{})
        if ok {
            var roomImageURLs []string
            for _, photo := range roomPhotos {
                photoMap, ok := photo.(map[string]interface{})
                if !ok {
                    continue
                }
                thumbURL, ok := photoMap["thumb_url"].(string)
                if !ok {
                    continue
                }
                roomImageURLs = append(roomImageURLs, thumbURL)
            }

            if len(roomImageURLs) > 0 {
                roomPropertyImage := models.PropertyImage{
                    PropertyID:   property.HotelID,
                    PropertyName: strconv.Itoa(property.HotelID),
                    ImageType:    "room_photos",
                    ImageURLs:    roomImageURLs,
                }
                allPropertyImages = append(allPropertyImages, roomPropertyImage)
            }
        }
    }

    log.Printf("Finished fetching images. Total images collected: %d", len(allPropertyImages))
    return allPropertyImages, nil
}