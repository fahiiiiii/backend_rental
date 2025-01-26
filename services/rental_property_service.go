package services

import (
	"encoding/json"
	"fmt"
	"os"
	"backend_rental/models"
)

type RentalPropertyService struct{}

func (s *RentalPropertyService) GenerateRentalPropertyJSON(properties []map[string]interface{}, propertyDetails []map[string]interface{}) error {
	rentalProperties := []models.RentalProperty{}

	for _, prop := range properties {
		for _, detail := range propertyDetails {
			if prop["id"] == detail["hotel_id"] {
				// Convert amenities to JSON string
				amenitiesJSON, _ := json.Marshal(convertToStringSlice(detail["amenities"]))
				
				rentalProp := models.RentalProperty{
					PropertyID:   int64(prop["id"].(float64)),
					Name:         prop["name"].(string),
					CityID:       prop["cityId"].(string),
					Bedrooms:     int(detail["bedrooms"].(float64)),
					Bathrooms:    int(detail["bathrooms"].(float64)),
					Amenities:    string(amenitiesJSON),
					PropertyType: detail["property_type"].(string),
				}
				rentalProperties = append(rentalProperties, rentalProp)
				break
			}
		}
	}

	return s.writeJSONFile(rentalProperties)
}

func (s *RentalPropertyService) writeJSONFile(data []models.RentalProperty) error {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile("data/RentalProperty.json", file, 0644)
	return err
}

func convertToStringSlice(data interface{}) []string {
	slice, ok := data.([]interface{})
	if !ok {
		return []string{}
	}
	
	result := make([]string, len(slice))
	for i, v := range slice {
		result[i] = v.(string)
	}
	return result
}

func convertToStringSgSlice(amenities interface{}) []string {
    switch v := amenities.(type) {
    case []interface{}:
        result := make([]string, len(v))
        for i, item := range v {
            result[i] = fmt.Sprintf("%v", item)
        }
        return result
    case []string:
        return v
    default:
        return []string{}
    }
}