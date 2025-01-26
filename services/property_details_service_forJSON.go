package services

import (
	"encoding/json"
	"os"
	"backend_rental/models"
	"backend_rental/utils"
)

type PropertyDetailsServiceJSON struct{}

func (s *PropertyDetailsServiceJSON) GeneratePropertyDetailsJSON(
	descImages []map[string]interface{}, 
	properties []map[string]interface{}, 
	propertyImages []map[string]interface{}) error {
	
	propertyDetailsList := []models.PropertyDetails{}

	for _, descImage := range descImages {
		propertyID := descImage["property_id"].(float64)

		// Find matching property details
		var matchedProperty map[string]interface{}
		for _, prop := range properties {
			if prop["id"] == propertyID {
				matchedProperty = prop
				break
			}
		}

		// Find matching property images
		var matchedImages []string
		var imageType string
		for _, img := range propertyImages {
			if img["property_id"] == propertyID {
				matchedImages = utils.ConvertToStringSlice(img["image_urls"])
				imageType = img["image_type"].(string)
				break
			}
		}

		propertyDetails := models.PropertyDetails{
			PropertyID:      int64(propertyID),
			Description:     descImage["description"].(string),
			ReviewScore:     matchedProperty["reviewScore"].(float64),
			ReviewCount:     int(matchedProperty["reviewCount"].(float64)),
			ReviewScoreWord: matchedProperty["reviewScoreWord"].(string),
			ImageType:       imageType,
			ImageUrls:       matchedImages,
		}

		propertyDetailsList = append(propertyDetailsList, propertyDetails)
	}

	return s.writeJSONFile(propertyDetailsList)
}


// Add this utility function
// func convertToStringSlice(data interface{}) []string {
// 	slice, ok := data.([]interface{})
// 	if !ok {
// 		return []string{}
// 	}
	
// 	result := make([]string, len(slice))
// 	for i, v := range slice {
// 		result[i] = v.(string)
// 	}
// 	return result
// }

func (s *PropertyDetailsServiceJSON) writeJSONFile(data []models.PropertyDetails) error {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile("data/PropertyDetails.json", file, 0644)
	return err
}
