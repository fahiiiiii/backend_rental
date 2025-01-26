package utils

import (
    "encoding/json"
    "os"
    "backend_rental/models"
)

func SavePropertyImagesToJSON(images []models.PropertyImage, filename string) error {
    file, err := json.MarshalIndent(images, "", " ")
    if err != nil {
        return err
    }
    return os.WriteFile(filename, file, 0644)
}

func LoadPropertiesFromJSON(filename string) ([]models.Property, error) {
    var properties []models.Property
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    err = json.Unmarshal(data, &properties)
    return properties, err
}