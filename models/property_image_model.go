package models

type PropertyImageResponse struct {
    Data map[string]interface{} `json:"data"`
}

type PropertyImage struct {
    PropertyID   int      `json:"property_id"`
    PropertyName string   `json:"property_name"`
    Type         string   `json:"type"`
    URLs         []string `json:"urls"`
}