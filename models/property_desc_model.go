package models

type PropertyDescription struct {
    PropertyID   int    `json:"property_id"`
    PropertyName string `json:"property_name"`
    Description  string `json:"description"`
}