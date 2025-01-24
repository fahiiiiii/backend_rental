package models

type PropertyImageResponse struct {
    Data map[string]interface{} `json:"data"`
    Meta struct {} `json:"meta"`
}

type PropertyImage struct {
    PropertyID   int      `json:"property_id"`
    PropertyName string   `json:"property_name"`
    Type         string   `json:"type"`
    Images       []string `json:"images"`
}
// package models

// type PropertyImageResponse struct {
//     Data map[string]interface{} `json:"data"`
// }

// type PropertyImage struct {
//     PropertyID   int      `json:"property_id"`
//     PropertyName string   `json:"property_name"`
//     ImageType    string   `json:"image_type"`
//     ImageURLs    []string `json:"image_urls"`
// }